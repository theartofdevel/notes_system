package note

import (
	"context"
)

type Storage interface {
	Create(ctx context.Context, note Note) (string, error)
	FindOne(ctx context.Context, uuid string) (Note, error)
	FindByCategoryUUID(ctx context.Context, uuid string) ([]Note, error)
	Update(ctx context.Context, note Note) error
	Delete(ctx context.Context, uuid string) error
}
