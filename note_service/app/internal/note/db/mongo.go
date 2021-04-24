package db

import (
	"context"
	"errors"
	"fmt"
	"github.com/theartofdevel/notes_system/note_service/internal/apperror"
	"github.com/theartofdevel/notes_system/note_service/internal/note"
	"github.com/theartofdevel/notes_system/note_service/pkg/logging"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"time"
)

var _ note.Storage = &db{}

type db struct {
	collection *mongo.Collection
	logger     logging.Logger
}

func NewStorage(storage *mongo.Database, collection string, logger logging.Logger) note.Storage {
	return &db{
		collection: storage.Collection(collection),
		logger:     logger,
	}
}

func (s *db) Create(ctx context.Context, note note.Note) (uuid string, err error) {
	nCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()
	result, err := s.collection.InsertOne(nCtx, note)
	if err != nil {
		return "", fmt.Errorf("failed to execute query. error: %w", err)
	}

	oid, ok := result.InsertedID.(primitive.ObjectID)
	if ok {
		return oid.Hex(), nil
	}
	return "", fmt.Errorf("failed to convet objectid to hex")
}

func (s *db) FindOne(ctx context.Context, uuid string) (n note.Note, err error) {
	objectID, err := primitive.ObjectIDFromHex(uuid)
	if err != nil {
		return n, fmt.Errorf("failed to convert hex to objectid. error: %w", err)
	}

	filter := bson.M{"_id": objectID}

	opts := options.FindOneOptions{
		Projection: bson.M{"short_body": 0},
	}

	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()
	result := s.collection.FindOne(ctx, filter, &opts)
	if result.Err() != nil {
		s.logger.Error(result.Err())
		if errors.Is(result.Err(), mongo.ErrNoDocuments) {
			return n, apperror.ErrNotFound
		}
		return n, fmt.Errorf("failed to execute query. error: %w", err)
	}
	if err = result.Decode(&n); err != nil {
		return n, fmt.Errorf("failed to decode document. error: %w", err)
	}

	return n, nil
}

func (s *db) FindByCategoryUUID(ctx context.Context, categoryUUID string) (notes []note.Note, err error) {
	opts := options.FindOptions{
		Projection: bson.M{"body": 0},
	}

	filter := bson.M{"category_uuid": bson.M{"$eq": categoryUUID}}

	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()
	cur, err := s.collection.Find(ctx, filter, &opts)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return notes, apperror.ErrNotFound
		}
		return notes, fmt.Errorf("failed to execute query. error: %w", err)
	}
	if err = cur.All(ctx, &notes); err == nil {
		return notes, nil
	}
	return notes, fmt.Errorf("failed to decode document. error: %w", err)
}

func (s *db) Update(ctx context.Context, note note.Note) error {
	objectID, err := primitive.ObjectIDFromHex(note.UUID)
	if err != nil {
		return fmt.Errorf("failed to parse note uuid due to error %w", err)
	}

	filter := bson.M{"_id": objectID}

	noteByte, err := bson.Marshal(note)
	if err != nil {
		return fmt.Errorf("failed to marshal document. error: %w", err)
	}

	var updateObj bson.M
	err = bson.Unmarshal(noteByte, &updateObj)
	if err != nil {
		return fmt.Errorf("failed to unmarshal document. error: %w", err)
	}

	delete(updateObj, "_id")

	update := bson.M{
		"$set": updateObj,
	}

	if note.Tags != nil {
		update["$set"].(bson.M)["tags"] = note.Tags
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
		return fmt.Errorf("failed to parse note uuid")
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