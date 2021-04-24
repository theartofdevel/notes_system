package file

import (
	"context"
	"github.com/theartofdevel/notes_system/file_service/pkg/logging"
)

var _ Service = &service{}

type service struct {
	storage Storage
	logger  logging.Logger
}

func NewService(noteStorage Storage, logger logging.Logger) (Service, error) {
	return &service{
		storage: noteStorage,
		logger:  logger,
	}, nil
}

type Service interface {
	GetFile(ctx context.Context, noteUUID, fileName string) (f *File, err error)
	GetFilesByNoteUUID(ctx context.Context, noteUUID string) ([]*File, error)
	Create(ctx context.Context, noteUUID string, dto CreateFileDTO) error
	Delete(ctx context.Context, noteUUID, fileName string) error
}

func (s *service) GetFile(ctx context.Context, noteUUID, fileId string) (f *File, err error) {
	f, err = s.storage.GetFile(ctx, noteUUID, fileId)
	if err != nil {
		return f, err
	}
	return f, nil
}

func (s *service) GetFilesByNoteUUID(ctx context.Context, noteUUID string) ([]*File, error) {
	files, err := s.storage.GetFilesByNoteUUID(ctx, noteUUID)
	if err != nil {
		return nil, err
	}
	return files, nil
}

func (s *service) Create(ctx context.Context, noteUUID string, dto CreateFileDTO) error {
	dto.NormalizeName()
	file, err := NewFile(dto)
	if err != nil {
		return err
	}
	err = s.storage.CreateFile(ctx, noteUUID, file)
	if err != nil {
		return err
	}
	return nil
}

func (s *service) Delete(ctx context.Context, noteUUID, fileName string) error {
	err := s.storage.DeleteFile(ctx, noteUUID, fileName)
	if err != nil {
		return err
	}
	return nil
}