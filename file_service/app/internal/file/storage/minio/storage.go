package minio

import (
	"bytes"
	"context"
	"fmt"
	"github.com/theartofdevel/notes_system/file_service/internal/apperror"
	"github.com/theartofdevel/notes_system/file_service/internal/file"
	"github.com/theartofdevel/notes_system/file_service/pkg/logging"
	"github.com/theartofdevel/notes_system/file_service/pkg/minio"
	"io"
)

type minioStorage struct {
	client *minio.Client
	logger logging.Logger
}

func NewStorage(endpoint, accessKeyID, secretAccessKey string, logger logging.Logger) (file.Storage, error) {
	client, err := minio.NewClient(endpoint, accessKeyID, secretAccessKey, logger)
	if err != nil {
		return nil, fmt.Errorf("failed to create minio client. err: %w", err)
	}
	return &minioStorage{
		client: client,
	}, nil
}

func (m *minioStorage) GetFile(ctx context.Context, bucketName, fileID string) (*file.File, error) {
	obj, err := m.client.GetFile(ctx, bucketName, fileID)
	if err != nil {
		return nil, fmt.Errorf("failed to get file. err: %w", err)
	}
	defer obj.Close()
	objectInfo, err := obj.Stat()
	if err != nil {
		return nil, fmt.Errorf("failed to get file. err: %w", err)
	}
	buffer := make([]byte, objectInfo.Size)
	_, err = obj.Read(buffer)
	if err != nil && err != io.EOF {
		return nil, fmt.Errorf("failed to get objects. err: %w", err)
	}
	f := file.File{
		ID:	   objectInfo.Key,
		Name:  objectInfo.UserMetadata["Name"],
		Size:  objectInfo.Size,
		Bytes: buffer,
	}
	return &f, nil
}

func (m *minioStorage) GetFilesByNoteUUID(ctx context.Context, noteUUID string) ([]*file.File, error) {
	objects, err := m.client.GetBucketFiles(ctx, noteUUID)
	if err != nil {
		return nil, fmt.Errorf("failed to get objects. err: %w", err)
	}
	if len(objects) == 0 {
		return nil, apperror.ErrNotFound
	}

	var files []*file.File
	for _, obj := range objects {
		stat, err := obj.Stat()
		if err != nil {
			m.logger.Errorf("failed to get objects. err: %v", err)
			continue
		}
		buffer := make([]byte, stat.Size)
		_, err = obj.Read(buffer)
		if err != nil && err != io.EOF {
			m.logger.Errorf("failed to get objects. err: %v", err)
			continue
		}
		f := file.File{
			ID:  stat.Key,
			Name: stat.UserMetadata["Name"],
			Size:  stat.Size,
			Bytes: buffer,
		}
		files = append(files, &f)
		obj.Close()
	}

	return files, nil
}

func (m *minioStorage) CreateFile(ctx context.Context, noteUUID string, file *file.File) error {
	err := m.client.UploadFile(ctx, file.ID, file.Name, noteUUID, file.Size, bytes.NewBuffer(file.Bytes))
	if err != nil {
		return err
	}
	return nil
}

func (m *minioStorage) DeleteFile(ctx context.Context, noteUUID, fileId string) error {
	err := m.client.DeleteFile(ctx, noteUUID, fileId)
	if err != nil {
		return err
	}
	return nil
}