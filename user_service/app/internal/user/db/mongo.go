package db

import (
	"context"
	"errors"
	"fmt"
	"github.com/theartofdevel/notes_system/user_service/internal/appError"
	"github.com/theartofdevel/notes_system/user_service/internal/user"
	"github.com/theartofdevel/notes_system/user_service/pkg/logging"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"time"
)

var _ user.Storage = &db{}

type db struct {
	collection *mongo.Collection
	client     *mongo.Client
	logger     logging.Logger
}

func NewStorage(ctx context.Context, hostname, port, username, password, authSource, database, entity string, logger logging.Logger) (user.Storage, error) {
	mongoDBURL := fmt.Sprintf("mongodb://%s:%s@%s:%s", username, password, hostname, port)
	credentials := options.Credential{
		AuthSource:  authSource,
		Username:    username,
		Password:    password,
		PasswordSet: true,
	}
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(mongoDBURL).SetAuth(credentials))
	if err != nil {
		return nil, fmt.Errorf("failed to create client to mongodb due to error %w", err)
	}

	collection := client.Database(database).Collection(entity)

	s := db{
		client:     client,
		collection: collection,
		logger:     logger,
	}

	err = s.isConnected(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to mongodb due to error %w", err)
	}
	return &s, nil
}

func (s *db) Create(ctx context.Context, dto user.CreateUserDTO) (string, error) {
	err := s.isConnected(ctx)
	if err != nil {
		return "", fmt.Errorf("storage is not connected. error: %w", err)
	}

	nCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()
	result, err := s.collection.InsertOne(nCtx, dto)
	if err != nil {
		return "", fmt.Errorf("failed to execute query. error: %w", err)
	}

	oid, ok := result.InsertedID.(primitive.ObjectID)
	if ok {
		return oid.Hex(), nil
	}
	return "", fmt.Errorf("failed to convet objectid to hex")
}

func (s *db) FindOne(ctx context.Context, uuid string) (u user.User, err error) {
	err = s.isConnected(ctx)
	if err != nil {
		return u, fmt.Errorf("storage is not connected. error: %w", err)
	}

	objectID, err := primitive.ObjectIDFromHex(uuid)
	if err != nil {
		return u, fmt.Errorf("failed to convert hex to objectid. error: %w", err)
	}

	filter := bson.M{"_id": objectID}

	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()
	result := s.collection.FindOne(ctx, filter)
	if result.Err() != nil {
		s.logger.Error(result.Err())
		if errors.Is(result.Err(), mongo.ErrNoDocuments) {
			return u, apperror.ErrNotFound
		}
		return u, fmt.Errorf("failed to execute query. error: %w", err)
	}
	if err = result.Decode(&u); err != nil {
		return u, fmt.Errorf("failed to decode document. error: %w", err)
	}

	return u, nil
}

func (s *db) FindByEmail(ctx context.Context, email string) (u user.User, err error) {
	err = s.isConnected(ctx)
	if err != nil {
		return u, fmt.Errorf("storage is not connected. error: %w", err)
	}

	filter := bson.M{"email": email}

	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	result := s.collection.FindOne(ctx, filter)
	err = result.Err()
	if err != nil {
		if errors.Is(result.Err(), mongo.ErrNoDocuments) {
			return u, apperror.ErrNotFound
		}
		return u, fmt.Errorf("failed to execute query. error: %w", err)
	}
	if err = result.Decode(&u); err != nil {
		return u, fmt.Errorf("failed to decode document. error: %w", err)
	}

	return u, nil
}

func (s *db) Update(ctx context.Context, uuid string, dto user.UpdateUserDTO) error {
	err := s.isConnected(ctx)
	if err != nil {
		return fmt.Errorf("storage is not connected")
	}

	objectID, err := primitive.ObjectIDFromHex(uuid)
	if err != nil {
		return fmt.Errorf("failed to execute query. error: %w", err)
	}

	filter := bson.M{"_id": objectID}

	userByte, err := bson.Marshal(dto)
	if err != nil {
		return fmt.Errorf("failed to marshal document. error: %w", err)
	}

	var updateObj bson.M
	err = bson.Unmarshal(userByte, &updateObj)
	if err != nil {
		return fmt.Errorf("failed to unmarshal document. error: %w", err)
	}

	update := bson.M{
		"$set": updateObj,
	}

	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()
	result, err := s.collection.UpdateOne(ctx, filter, update)
	if err != nil {
		return fmt.Errorf("failed to execute query. error: %w", err)
	}
	if result.MatchedCount == 0 {
		return apperror.ErrNotFound
	}

	s.logger.Tracef("Matched %v documents and updated %v documents.\n", result.MatchedCount, result.ModifiedCount)

	return nil
}

func (s *db) Delete(ctx context.Context, uuid string) error {
	err := s.isConnected(ctx)
	if err != nil {
		return fmt.Errorf("storage is not connected. error: %w", err)
	}

	objectID, err := primitive.ObjectIDFromHex(uuid)
	if err != nil {
		return fmt.Errorf("failed to convet objectid to hex. error: %w", err)
	}
	filter := bson.M{"_id": objectID}

	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()
	result, err := s.collection.DeleteOne(ctx, filter)
	if err != nil {
		return fmt.Errorf("failed to execute query")
	}
	if result.DeletedCount == 0 {
		return apperror.ErrNotFound
	}

	s.logger.Tracef("Delete %v documents.\n", result.DeletedCount)

	return nil
}

func (s *db) isConnected(ctx context.Context) error {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()
	return s.client.Ping(ctx, nil)
}

func (s *db) Close(ctx context.Context) error {
	var err error
	if err = s.isConnected(ctx); err != nil {
		ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
		defer cancel()
		err = s.client.Disconnect(ctx)
	}
	s.collection = nil
	s.client = nil
	return err
}
