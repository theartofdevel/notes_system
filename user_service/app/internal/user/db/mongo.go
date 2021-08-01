package db

import (
	"context"
	"errors"
	"fmt"
	"github.com/theartofdevel/notes_system/user_service/internal/apperror"
	"github.com/theartofdevel/notes_system/user_service/internal/user"
	"github.com/theartofdevel/notes_system/user_service/pkg/logging"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"time"
)

var _ user.Storage = &db{}

type db struct {
	collection *mongo.Collection
	logger     logging.Logger
}

func NewStorage(storage *mongo.Database, collection string, logger logging.Logger) user.Storage {
	return &db{
		collection: storage.Collection(collection),
		logger:     logger,
	}
}

func (s *db) Create(ctx context.Context, user user.User) (string, error) {
	nCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()
	result, err := s.collection.InsertOne(nCtx, user)
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

func (s *db) Update(ctx context.Context, user user.User) error {
	objectID, err := primitive.ObjectIDFromHex(user.UUID)
	if err != nil {
		return fmt.Errorf("failed to execute query. error: %w", err)
	}

	filter := bson.M{"_id": objectID}
	userByte, err := bson.Marshal(user)
	if err != nil {
		return fmt.Errorf("failed to marshal document. error: %w", err)
	}

	var updateObj bson.M
	err = bson.Unmarshal(userByte, &updateObj)
	if err != nil {
		return fmt.Errorf("failed to unmarshal document. error: %w", err)
	}

	delete(updateObj, "_id")

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