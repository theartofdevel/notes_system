package file

import (
	"context"
)

type Storage interface {
	GetFile(ctx context.Context, bucketName, fileName string) (*File, error)
	GetFilesByNoteUUID(ctx context.Context, uuid string) ([]*File, error)
	CreateFile(ctx context.Context, noteUUID string, file *File) error
	DeleteFile(ctx context.Context, noteUUID, fileName string) error
}