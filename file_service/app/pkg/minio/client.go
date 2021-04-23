package minio

import (
	"context"
	"fmt"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"github.com/theartofdevel/notes_system/file_service/pkg/logging"
	"io"
	"time"
)

type Object struct {
	ID   string
	Size int64
	Tags map[string]string
}

type Client struct {
	logger      logging.Logger
	minioClient *minio.Client
}

func NewClient(endpoint, accessKeyID, secretAccessKey string, logger logging.Logger) (*Client, error) {
	minioClient, err := minio.New(endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(accessKeyID, secretAccessKey, ""),
		Secure: false,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create minio client. err: %w", err)
	}

	return &Client{
		logger:      logger,
		minioClient: minioClient,
	}, nil
}

func (c *Client) GetFile(ctx context.Context, bucketName, fileId string) (*minio.Object, error) {
	reqCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	obj, err := c.minioClient.GetObject(reqCtx, bucketName, fileId, minio.GetObjectOptions{})
	if err != nil {
		return nil, fmt.Errorf("failed to get file with id: %s from minio bucket %s. err: %w", fileId, bucketName, err)
	}
	return obj, nil
}

func (c *Client) GetBucketFiles(ctx context.Context, bucketName string) ([]*minio.Object, error) {
	reqCtx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	var files []*minio.Object
	for lobj := range c.minioClient.ListObjects(reqCtx, bucketName, minio.ListObjectsOptions{WithMetadata: true}) {
		if lobj.Err != nil {
			c.logger.Errorf("failed to list object from minio bucket %s. err: %v", bucketName, lobj.Err)
			continue
		}
		object, err := c.minioClient.GetObject(ctx, bucketName, lobj.Key, minio.GetObjectOptions{})
		if err != nil {
			c.logger.Errorf("failed to get object key=%s from minio bucket %s. err: %v", lobj.Key, bucketName, lobj.Err)
			continue
		}
		files = append(files, object)
	}
	return files, nil
}

func (c *Client) UploadFile(ctx context.Context, fileId, fileName, bucketName string, fileSize int64, reader io.Reader) error {
	reqCtx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	exists, errBucketExists := c.minioClient.BucketExists(ctx, bucketName)
	if errBucketExists != nil || !exists {
		c.logger.Warnf("no bucket %s. creating new one...", bucketName)
		err := c.minioClient.MakeBucket(ctx, bucketName, minio.MakeBucketOptions{})
		if err != nil {
			return fmt.Errorf("failed to create new bucket. err: %w", err)
		}
	}

	c.logger.Debugf("put new object %s to bucket %s", fileName, bucketName)
	_, err := c.minioClient.PutObject(reqCtx, bucketName, fileId, reader, fileSize,
		minio.PutObjectOptions{
			UserMetadata: map[string]string{
				"Name": fileName,
			},
			ContentType: "application/octet-stream",
		})
	if err != nil {
		return fmt.Errorf("failed to upload file. err: %w", err)
	}
	return nil
}

func (c *Client) DeleteFile(ctx context.Context, noteUUID, fileName string) error {
	err := c.minioClient.RemoveObject(ctx, noteUUID, fileName, minio.RemoveObjectOptions{})
	if err != nil {
		return fmt.Errorf("failed to delete file. err: %w", err)
	}
	return nil
}
