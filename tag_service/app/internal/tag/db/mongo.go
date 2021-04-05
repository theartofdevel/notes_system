package db

import (
	"context"
	"errors"
	"fmt"
	"github.com/theartofdevel/notes_system/tag_service/internal/apperror"
	"github.com/theartofdevel/notes_system/tag_service/internal/tag"
	"github.com/theartofdevel/notes_system/tag_service/pkg/logging"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"time"
)

var _ tag.Storage = &db{}

type db struct {
	collection *mongo.Collection
	client     *mongo.Client
	logger     logging.Logger
}

func NewStorage(ctx context.Context, hostname, port, username, password, authSource, database, entity string, logger logging.Logger) (tag.Storage, error) {
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

func (s *db) Create(ctx context.Context, dto tag.CreateTagDTO) (id int, err error) {
	err = s.isConnected(ctx)
	if err != nil {
		return id, fmt.Errorf("storage is not connected. error: %w", err)
	}

	nCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	findOptions := options.FindOptions{}
	findOptions.SetSort(bson.D{{"_id", -1}}) //nolint:govet
	findOptions.SetLimit(1)
	var nTag tag.Tag
	cursor, err := s.collection.Find(nCtx, bson.M{}, &findOptions)
	if err != nil {
		return id, fmt.Errorf("failed to execute query. error: %w", err)
	}

	nTagID := 1
	tryCount := 3
	for tryCount >= 0 {
		if cursor.Next(ctx) {
			err = cursor.Decode(&nTag)
			if err != nil {
				return 0, err
			}
			nTagID = nTag.ID + 1
		} else if tryCount < 3 {
			return id, fmt.Errorf("duplicate key error")
		}

		tryCount--

		dto.ID = nTagID
		nCtx, cancel = context.WithTimeout(ctx, 5*time.Second)
		defer cancel()
		_, err = s.collection.InsertOne(nCtx, dto)
		if err != nil {
			if mongo.IsDuplicateKeyError(err) {
				s.logger.Warnf("duplicate key error. continue optimistic loop")
				continue
			} else {
				return id, fmt.Errorf("failed to execute query. error: %w", err)
			}
		}
		break
	}

	return nTagID, nil
}

func (s *db) FindOne(ctx context.Context, id int) (t tag.Tag, err error) {
	err = s.isConnected(ctx)
	if err != nil {
		return t, fmt.Errorf("storage is not connected. error: %w", err)
	}

	filter := bson.M{"_id": id}

	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	result := s.collection.FindOne(ctx, filter)
	if result.Err() != nil {
		s.logger.Error(result.Err())
		if errors.Is(result.Err(), mongo.ErrNoDocuments) {
			return t, apperror.ErrNotFound
		}
		return t, fmt.Errorf("failed to execute query. error: %w", err)
	}
	if err := result.Decode(&t); err != nil {
		return t, fmt.Errorf("failed to decode document. error: %w", err)
	}

	return t, nil
}

func (s *db) FindMany(ctx context.Context, ids []int) (tags []tag.Tag, err error) {
	err = s.isConnected(ctx)
	if err != nil {
		return tags, fmt.Errorf("storage is not connected. error: %w", err)
	}

	filter := bson.M{"_id": bson.M{"$in": ids}}

	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	cur, err := s.collection.Find(ctx, filter)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return tags, apperror.ErrNotFound
		}
		return tags, fmt.Errorf("failed to execute query. error: %w", err)
	}
	if err = cur.All(ctx, &tags); err == nil {
		return tags, nil
	}
	return tags, fmt.Errorf("failed to decode document. error: %w", err)
}

func (s *db) Update(ctx context.Context, id int, dto tag.UpdateTagDTO) error {
	err := s.isConnected(ctx)
	if err != nil {
		return fmt.Errorf("storage is not connected. error: %w", err)
	}

	filter := bson.M{"_id": id}

	tagByte, err := bson.Marshal(dto)
	if err != nil {
		return fmt.Errorf("failed to marshal document. error: %w", err)
	}

	var updateObj bson.M
	err = bson.Unmarshal(tagByte, &updateObj)
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

func (s *db) Delete(ctx context.Context, id int) error {
	err := s.isConnected(ctx)
	if err != nil {
		return fmt.Errorf("storage is not connected. error: %w", err)
	}

	filter := bson.M{"_id": id}

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
