package tag

import (
	"context"
)

type Storage interface {
	Create(ctx context.Context, dto CreateTagDTO) (int, error)
	FindOne(ctx context.Context, id int) (Tag, error)
	FindMany(ctx context.Context, ids []int) ([]Tag, error)
	Update(ctx context.Context, id int, dto UpdateTagDTO) error
	Delete(ctx context.Context, id int) error
}
