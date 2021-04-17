package user_service

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

var _ UserService = &client{}

type client struct {
	base     rest.BaseClient
	Resource string
}

func NewService(baseURL string, resource string, logger logging.Logger) UserService {
	c := client{
		Resource: resource,
		base: rest.BaseClient{
			BaseURL: baseURL,
			HTTPClient: &http.Client{
				Timeout: 10 * time.Second,
			},
			Logger: logger,
		},
	}
	return &c
}

type UserService interface {
	GetByEmailAndPassword(ctx context.Context, email, password string) (User, error)
	GetByUUID(ctx context.Context, uuid string) (User, error)
	Create(ctx context.Context, dto CreateUserDTO) (User, error)
	Update(ctx context.Context, uuid string, dto UpdateUserDTO) error
	Delete(ctx context.Context, uuid string) error
}

func (c *client) GetByEmailAndPassword(ctx context.Context, email, password string) (u User, err error) {
	c.base.Logger.Debug("add email and password to filter options")
	filters := []rest.FilterOptions{
		{
			Field:  "email",
			Values: []string{email},
		},
		{
			Field:  "password",
			Values: []string{password},
		},
	}

	c.base.Logger.Debug("build url with resource and filter")
	uri, err := c.base.BuildURL(c.Resource, filters)
	if err != nil {
		return u, fmt.Errorf("failed to build URL. error: %v", err)
	}
	c.base.Logger.Tracef("url: %s", uri)

	c.base.Logger.Debug("create new request")
	req, err := http.NewRequest(http.MethodGet, uri, nil)
	if err != nil {
		return u, fmt.Errorf("failed to create new request due to error: %w", err)
	}

	c.base.Logger.Debug("send request")
	reqCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()
	req = req.WithContext(reqCtx)
	response, err := c.base.SendRequest(req)
	if err != nil {
		return u, fmt.Errorf("failed to send request due to error: %w", err)
	}

	if response.IsOk {
		if err = json.NewDecoder(response.Body()).Decode(&u); err != nil {
			return u, fmt.Errorf("failed to decode body due to error %w", err)
		}
		return u, nil
	}
	return u, apperror.APIError(response.Error.ErrorCode, response.Error.Message, response.Error.DeveloperMessage)
}

func (c *client) GetByUUID(ctx context.Context, uuid string) (User, error) {
	var u User

	c.base.Logger.Debug("build url with resource and filter")
	uri, err := c.base.BuildURL(fmt.Sprintf("%s/%s", c.Resource, uuid), nil)
	if err != nil {
		return u, fmt.Errorf("failed to build URL. error: %v", err)
	}
	c.base.Logger.Tracef("url: %s", uri)

	c.base.Logger.Debug("create new request")
	req, err := http.NewRequest("GET", uri, nil)
	if err != nil {
		return u, fmt.Errorf("failed to create new request due to error: %w", err)
	}

	c.base.Logger.Debug("send request")
	reqCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()
	req = req.WithContext(reqCtx)
	response, err := c.base.SendRequest(req)
	if err != nil {
		return u, fmt.Errorf("failed to send request due to error: %w", err)
	}

	if response.IsOk {
		defer response.Body().Close()
		if err = json.NewDecoder(response.Body()).Decode(&u); err != nil {
			return u, fmt.Errorf("failed to decode body due to error %w", err)
		}
		return u, nil
	}
	return u, apperror.APIError(response.Error.ErrorCode, response.Error.Message, response.Error.DeveloperMessage)
}

func (c *client) Create(ctx context.Context, dto CreateUserDTO) (User, error) {
	var u User

	c.base.Logger.Debug("build url with resource and filter")
	uri, err := c.base.BuildURL(c.Resource, nil)
	if err != nil {
		return u, fmt.Errorf("failed to build URL. error: %v", err)
	}
	c.base.Logger.Tracef("url: %s", uri)

	c.base.Logger.Debug("convert dto to map")
	structs.DefaultTagName = "json"
	data := structs.Map(dto)

	c.base.Logger.Debug("marshal map to bytes")
	dataBytes, err := json.Marshal(data)
	if err != nil {

		return u, fmt.Errorf("failed to marshal dto")
	}

	c.base.Logger.Debug("create new request")
	req, err := http.NewRequest(http.MethodPost, uri, bytes.NewBuffer(dataBytes))
	if err != nil {
		return u, fmt.Errorf("failed to create new request due to error: %w", err)
	}

	c.base.Logger.Debug("send request")
	reqCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()
	req = req.WithContext(reqCtx)
	response, err := c.base.SendRequest(req)
	if err != nil {
		return u, fmt.Errorf("failed to send request due to error: %w", err)
	}

	if response.IsOk {
		c.base.Logger.Debug("parse location header")
		userURL, err := response.Location()
		if err != nil {
			return u, fmt.Errorf("failed to get Location header")
		}
		c.base.Logger.Tracef("Location: %s", userURL.String())

		splitCategoryURL := strings.Split(userURL.String(), "/")
		userUUID := splitCategoryURL[len(splitCategoryURL)-1]
		u, err = c.GetByUUID(ctx, userUUID)
		if err != nil {
			return u, err
		}
		return u, nil
	}
	return u, apperror.APIError(response.Error.ErrorCode, response.Error.Message, response.Error.DeveloperMessage)
}

func (c *client) Update(ctx context.Context, uuid string, dto UpdateUserDTO) error {
	c.base.Logger.Debug("build url with resource and filter")
	uri, err := c.base.BuildURL(fmt.Sprintf("%s/%s", c.Resource, uuid), nil)
	if err != nil {
		return fmt.Errorf("failed to build URL. error: %v", err)
	}
	c.base.Logger.Tracef("url: %s", uri)

	c.base.Logger.Debug("convert dto to map")
	structs.DefaultTagName = "json"
	data := structs.Map(dto)

	c.base.Logger.Debug("marshal map to bytes")
	dataBytes, err := json.Marshal(data)
	if err != nil {
		return fmt.Errorf("failed to marshal dto")
	}

	c.base.Logger.Debug("create new request")
	req, err := http.NewRequest(http.MethodPatch, uri, bytes.NewBuffer(dataBytes))
	if err != nil {
		return fmt.Errorf("failed to create new request due to error: %w", err)
	}

	c.base.Logger.Debug("send request")
	reqCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()
	req = req.WithContext(reqCtx)
	response, err := c.base.SendRequest(req)
	if err != nil {
		return fmt.Errorf("failed to send request due to error: %w", err)
	}
	if response.IsOk {
		return nil
	}
	return apperror.APIError(response.Error.ErrorCode, response.Error.Message, response.Error.DeveloperMessage)
}

func (c *client) Delete(ctx context.Context, uuid string) error {
	c.base.Logger.Debug("build url with resource and filter")
	uri, err := c.base.BuildURL(fmt.Sprintf("%s/%s", c.Resource, uuid), nil)
	if err != nil {
		return fmt.Errorf("failed to build URL. error: %v", err)
	}
	c.base.Logger.Tracef("url: %s", uri)

	c.base.Logger.Debug("create new request")
	req, err := http.NewRequest("DELETE", uri, nil)
	if err != nil {
		return fmt.Errorf("failed to create new request due to error: %w", err)
	}

	c.base.Logger.Debug("send request")
	reqCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()
	req = req.WithContext(reqCtx)
	response, err := c.base.SendRequest(req)
	if err != nil {
		return fmt.Errorf("failed to send request due to error: %w", err)
	}

	if response.IsOk {
		return nil
	}
	return apperror.APIError(response.Error.ErrorCode, response.Error.Message, response.Error.DeveloperMessage)
}
