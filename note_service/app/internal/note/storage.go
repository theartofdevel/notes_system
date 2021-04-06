package note

import (
	"context"
)

type Storage interface {
	Create(ctx context.Context, dto CreateNoteDTO) (string, error)
	FindOne(ctx context.Context, uuid string) (Note, error)
	FindByCategoryUUID(ctx context.Context, uuid string) ([]Note, error)
	Update(ctx context.Context, uuid string, dto UpdateNoteDTO, tagsUpdate bool) error
	Delete(ctx context.Context, uuid string) error
}
