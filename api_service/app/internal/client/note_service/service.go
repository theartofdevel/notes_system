package note_service

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

var _ NoteService = &client{}

type client struct {
	Resource string
	base     rest.BaseClient
}

func NewService(baseURL string, resource string, logger logging.Logger) NoteService {
	return &client{
		Resource: resource,
		base: rest.BaseClient{
			BaseURL: baseURL,
			HTTPClient: &http.Client{
				Timeout: 10 * time.Second,
			},
			Logger: logger,
		},
	}
}

type NoteService interface {
	GetByCategoryUUID(ctx context.Context, categoryUUID string) ([]byte, error)
	GetByUUID(ctx context.Context, uuid string) ([]byte, error)
	Create(ctx context.Context, note CreateNoteDTO) (string, error)
	Update(ctx context.Context, uuid string, note UpdateNoteDTO) error
	Delete(ctx context.Context, uuid string) error
}

func (c *client) GetByCategoryUUID(ctx context.Context, categoryUUID string) ([]byte, error) {
	var notes []byte

	c.base.Logger.Debug("add email and password to filter options")
	filters := []rest.FilterOptions{
		{
			Field:  "category_uuid",
			Values: []string{categoryUUID},
		},
	}

	c.base.Logger.Debug("build url with resource and filter")
	uri, err := c.base.BuildURL(c.Resource, filters)
	if err != nil {
		return notes, fmt.Errorf("failed to build URL. error: %v", err)
	}
	c.base.Logger.Tracef("url: %s", uri)

	c.base.Logger.Debug("create new request")
	req, err := http.NewRequest("GET", uri, nil)
	if err != nil {
		return notes, fmt.Errorf("failed to create new request due to error: %v", err)
	}

	c.base.Logger.Debug("send request")
	reqCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()
	req = req.WithContext(reqCtx)
	response, err := c.base.SendRequest(req)
	if err != nil {
		return notes, fmt.Errorf("failed to send request due to error: %v", err)
	}

	if response.IsOk {
		c.base.Logger.Debug("read body")
		notes, err = response.ReadBody()
		if err != nil {
			return nil, fmt.Errorf("failed to read body")
		}
		return notes, nil
	}
	return nil, apperror.APIError(response.Error.ErrorCode, response.Error.Message, response.Error.DeveloperMessage)
}

func (c *client) GetByUUID(ctx context.Context, uuid string) ([]byte, error) {
	var note []byte

	c.base.Logger.Debug("build url with resource and filter")
	uri, err := c.base.BuildURL(fmt.Sprintf("%s/%s", c.Resource, uuid), nil)
	if err != nil {
		return note, fmt.Errorf("failed to build URL. error: %v", err)
	}
	c.base.Logger.Tracef("url: %s", uri)

	c.base.Logger.Debug("create new request")
	req, err := http.NewRequest("GET", uri, nil)
	if err != nil {
		return note, fmt.Errorf("failed to create new request due to error: %v", err)
	}

	c.base.Logger.Debug("send request")
	reqCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()
	req = req.WithContext(reqCtx)
	response, err := c.base.SendRequest(req)
	if err != nil {
		return note, fmt.Errorf("failed to send request due to error: %v", err)
	}

	if response.IsOk {
		c.base.Logger.Debug("read body")
		note, err = response.ReadBody()
		if err != nil {
			return nil, fmt.Errorf("failed to read body")
		}
		return note, nil
	}
	return nil, apperror.APIError(response.Error.ErrorCode, response.Error.Message, response.Error.DeveloperMessage)
}

func (c *client) Create(ctx context.Context, note CreateNoteDTO) (string, error) {
	var noteUUID string

	c.base.Logger.Debug("build url with resource and filter")
	uri, err := c.base.BuildURL(c.Resource, nil)
	if err != nil {
		return noteUUID, fmt.Errorf("failed to build URL. error: %v", err)
	}
	c.base.Logger.Tracef("url: %s", uri)

	c.base.Logger.Debug("convert dto to map")
	structs.DefaultTagName = "json"
	data := structs.Map(note)

	c.base.Logger.Debug("marshal map to bytes")
	dataBytes, err := json.Marshal(data)
	if err != nil {
		return noteUUID, fmt.Errorf("failed to marshal dto")
	}

	c.base.Logger.Debug("create new request")
	req, err := http.NewRequest("POST", uri, bytes.NewBuffer(dataBytes))
	if err != nil {
		return noteUUID, fmt.Errorf("failed to create new request due to error: %v", err)
	}

	c.base.Logger.Debug("send request")
	reqCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()
	req = req.WithContext(reqCtx)
	response, err := c.base.SendRequest(req)
	if err != nil {
		return noteUUID, fmt.Errorf("failed to send request due to error: %v", err)
	}

	if response.IsOk {
		c.base.Logger.Debug("parse location header")
		noteURL, err := response.Location()
		if err != nil {
			return noteUUID, fmt.Errorf("failed to get Location header")
		}
		c.base.Logger.Tracef("Location: %s", noteURL.String())

		splitCategoryURL := strings.Split(noteURL.String(), "/")
		noteUUID = splitCategoryURL[len(splitCategoryURL)-1]
		return noteUUID, nil
	}
	return noteUUID, apperror.APIError(response.Error.ErrorCode, response.Error.Message, response.Error.DeveloperMessage)
}

func (c *client) Update(ctx context.Context, uuid string, note UpdateNoteDTO) error {
	c.base.Logger.Debug("build url with resource and filter")
	uri, err := c.base.BuildURL(fmt.Sprintf("%s/%s", c.Resource, uuid), nil)
	if err != nil {
		return fmt.Errorf("failed to build URL. error: %v", err)
	}
	c.base.Logger.Tracef("url: %s", uri)

	c.base.Logger.Debug("convert dto to map")
	structs.DefaultTagName = "json"
	data := structs.Map(note)

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
