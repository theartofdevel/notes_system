package tag

import (
	"context"
	"errors"
	"fmt"
	"github.com/theartofdevel/notes_system/tag_service/internal/apperror"
	"github.com/theartofdevel/notes_system/tag_service/pkg/logging"
)

var _ Service = &service{}

type service struct {
	storage Storage
	logger  logging.Logger
}

func NewService(tagStorage Storage, logger logging.Logger) (Service, error) {
	return &service{
		storage: tagStorage,
		logger:  logger,
	}, nil
}

type Service interface {
	Create(ctx context.Context, dto CreateTagDTO) (int, error)
	GetOne(ctx context.Context, id int) (Tag, error)
	GetMany(ctx context.Context, ids []int) ([]Tag, error)
	Update(ctx context.Context, dto UpdateTagDTO) error
	Delete(ctx context.Context, id int) error
}

func (s service) Create(ctx context.Context, dto CreateTagDTO) (tagID int, err error) {
	tag := NewTag(dto)

	tagID, err = s.storage.Create(ctx, tag)

	if err != nil {
		if errors.Is(err, apperror.ErrNotFound) {
			return tagID, err
		}
		return tagID, fmt.Errorf("failed to create tag. error: %w", err)
	}

	return tagID, nil
}

func (s service) GetOne(ctx context.Context, id int) (t Tag, err error) {
	t, err = s.storage.FindOne(ctx, id)

	if err != nil {
		if errors.Is(err, apperror.ErrNotFound) {
			return t, err
		}
		return t, fmt.Errorf("failed to get one tag by id. error: %w", err)
	}
	return t, nil
}

func (s service) GetMany(ctx context.Context, ids []int) (tags []Tag, err error) {
	tags, err = s.storage.FindMany(ctx, ids)

	if err != nil {
		if errors.Is(err, apperror.ErrNotFound) {
			return tags, err
		}
		return tags, fmt.Errorf("failed to get many tags by ids. error: %w", err)
	}
	if len(tags) == 0 {
		return tags, apperror.ErrNotFound
	}

	return tags, nil
}

func (s service) Update(ctx context.Context, dto UpdateTagDTO) error {
	if dto.Name == "" && dto.Color == "" {
		return apperror.BadRequestError("no data to update")
	}

	tag := UpdatedTag(dto)

	err := s.storage.Update(ctx, tag)

	if err != nil {
		if errors.Is(err, apperror.ErrNotFound) {
			return err
		}
		return fmt.Errorf("failed to update tag. error: %w", err)
	}
	return nil
}

func (s service) Delete(ctx context.Context, id int) error {
	err := s.storage.Delete(ctx, id)

	if err != nil {
		if errors.Is(err, apperror.ErrNotFound) {
			return err
		}
		return fmt.Errorf("failed to delete tag. error: %w", err)
	}
	return nil
}
