package tag

import (
	"context"
)

type Storage interface {
	Create(ctx context.Context, t Tag) (int, error)
	FindOne(ctx context.Context, id int) (Tag, error)
	FindMany(ctx context.Context, ids []int) ([]Tag, error)
	Update(ctx context.Context, t Tag) error
	Delete(ctx context.Context, id int) error
}
