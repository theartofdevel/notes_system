package user

import (
	"context"
)

type Storage interface {
	Create(ctx context.Context, user CreateUserDTO) (string, error)
	FindByEmail(ctx context.Context, email string) (User, error)
	FindOne(ctx context.Context, uuid string) (User, error)
	Update(ctx context.Context, uuid string, user UpdateUserDTO) error
	Delete(ctx context.Context, uuid string) error
}
