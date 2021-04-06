package category_service

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/fatih/structs"
	"github.com/theartofdevel/notes_system/api_service/internal/apperror"
	"github.com/theartofdevel/notes_system/api_service/pkg/logging"
	"github.com/theartofdevel/notes_system/api_service/pkg/rest"
	"net/http"
	"strings"
	"time"
)

type client struct {
	base rest.BaseClient
	Resource string
}

func NewService(baseURL string, resource string, logger logging.Logger) CategoryService {
	return &client{
		Resource: resource,
		base: rest.BaseClient{
			BaseURL:  baseURL,
			HTTPClient: &http.Client{
				Timeout: 10 * time.Second,
			},
			Logger: logger,
		},
	}
}

type CategoryService interface {
	GetUserCategories(ctx context.Context, userUuid string) ([]byte, error)
	CreateCategory(ctx context.Context, dto CreateCategoryDTO) (string, error)
	UpdateCategory(ctx context.Context, uuid string, dto UpdateCategoryDTO) error
	DeleteCategory(ctx context.Context, dto DeleteCategoryDTO) error
}

func (c *client) GetUserCategories(ctx context.Context, userUuid string) ([]byte, error) {
	var categories []byte

	c.base.Logger.Debug("add user_uuid to filter options")
	filters := []rest.FilterOptions{
		{
			Field:  "user_uuid",
			Values: []string{userUuid},
		},
	}

	c.base.Logger.Debug("build url with resource and filter")
	uri, err := c.base.BuildURL(c.Resource, filters)
	if err != nil {
		return categories, fmt.Errorf("failed to build URL. error: %v", err)
	}
	c.base.Logger.Tracef("url: %s", uri)

	c.base.Logger.Debug("create new request")
	req, err := http.NewRequest("GET", uri, nil)
	if err != nil {
		return categories, fmt.Errorf("failed to create new request due to error: %v", err)
	}

	c.base.Logger.Debug("send request")
	reqCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()
	req = req.WithContext(reqCtx)
	response, err := c.base.SendRequest(req)
	if err != nil {
		return categories, fmt.Errorf("failed to send request due to error: %v", err)
	}

	if response.IsOk {
		c.base.Logger.Debug("read body")
		categories, err = response.ReadBody()
		if err != nil {
			return nil, fmt.Errorf("failed to read body")
		}
		return categories, nil
	}
	return nil, apperror.APIError(response.Error.ErrorCode, response.Error.Message, response.Error.DeveloperMessage)
}

func (c *client) CreateCategory(ctx context.Context, dto CreateCategoryDTO) (string, error) {
	var categoryUuid string

	c.base.Logger.Debug("build url with resource and filter")
	uri, err := c.base.BuildURL(c.Resource, nil)
	if err != nil {
		return categoryUuid, fmt.Errorf("failed to build URL. error: %v", err)
	}

	c.base.Logger.Debug("convert dto to map")
	structs.DefaultTagName = "json"
	data := structs.Map(dto)

	c.base.Logger.Debug("marshal map to bytes")
	dataBytes, err := json.Marshal(data)
	if err != nil {
		return categoryUuid, fmt.Errorf("failed to marshal dto")
	}

	c.base.Logger.Debug("create new request")
	req, err := http.NewRequest("POST", uri, bytes.NewBuffer(dataBytes))
	if err != nil {
		return categoryUuid, fmt.Errorf("failed to create new request due to error: %v", err)
	}

	c.base.Logger.Debug("send request")
	reqCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()
	req = req.WithContext(reqCtx)
	response, err := c.base.SendRequest(req)
	if err != nil {
		return categoryUuid, fmt.Errorf("failed to send request due to error: %w", err)
	}

	if response.IsOk {
		c.base.Logger.Debug("parse location header")
		categoryURL, err := response.Location()
		if err != nil {
			return categoryUuid, fmt.Errorf("failed to get Location header")
		}
		c.base.Logger.Tracef("Location: %s", categoryURL.String())

		splitCategoryURL := strings.Split(categoryURL.String(), "/")
		categoryUuid = splitCategoryURL[len(splitCategoryURL)-1]
		return categoryUuid, nil
	}
	return categoryUuid, apperror.APIError(response.Error.ErrorCode, response.Error.Message, response.Error.DeveloperMessage)
}

func (c *client) UpdateCategory(ctx context.Context, uuid string, dto UpdateCategoryDTO) error {
	c.base.Logger.Debug("build url with resource and filter")
	uri, err := c.base.BuildURL(fmt.Sprintf("%s/%s", c.Resource, uuid), nil)
	if err != nil {
		return fmt.Errorf("failed to build URL. error: %v", err)
	}

	c.base.Logger.Debug("convert dto to map")
	structs.DefaultTagName = "json"
	data := structs.Map(dto)

	c.base.Logger.Debug("marshal map to bytes")
	dataBytes, err := json.Marshal(data)
	if err != nil {
		return fmt.Errorf("failed to marshal dto")
	}

	c.base.Logger.Debug("create new request")
	req, err := http.NewRequest("PATCH", uri, bytes.NewBuffer(dataBytes))
	if err != nil {
		return fmt.Errorf("failed to create new request due to error: %v", err)
	}

	c.base.Logger.Debug("send request")
	reqCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()
	req = req.WithContext(reqCtx)
	response, err := c.base.SendRequest(req)
	if err != nil {
		return fmt.Errorf("failed to send request due to error: %v", err)
	}

	if response.IsOk {
		return nil
	}
	return apperror.APIError(response.Error.ErrorCode, response.Error.Message, response.Error.DeveloperMessage)
}

func (c *client) DeleteCategory(ctx context.Context, dto DeleteCategoryDTO) error {
	c.base.Logger.Debug("build url with resource and filter")
	uri, err := c.base.BuildURL(fmt.Sprintf("%s/%s", c.Resource, dto.Uuid), nil)
	if err != nil {
		return fmt.Errorf("failed to build URL. error: %v", err)
	}

	c.base.Logger.Debug("convert dto to map")
	structs.DefaultTagName = "json"
	data := structs.Map(dto)

	c.base.Logger.Debug("marshal map to bytes")
	dataBytes, err := json.Marshal(data)
	if err != nil {
		return fmt.Errorf("failed to marshal dto")
	}

	c.base.Logger.Debug("create new request")
	req, err := http.NewRequest("DELETE", uri, bytes.NewBuffer(dataBytes))
	if err != nil {
		return fmt.Errorf("failed to create new request due to error: %v", err)
	}

	c.base.Logger.Debug("send request")
	reqCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()
	req = req.WithContext(reqCtx)
	response, err := c.base.SendRequest(req)
	if err != nil {
		return fmt.Errorf("failed to send request due to error: %v", err)
	}

	if response.IsOk {
		return nil
	}
	return apperror.APIError(response.Error.ErrorCode, response.Error.Message, response.Error.DeveloperMessage)
}
