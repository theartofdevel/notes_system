package category

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/fatih/structs"
	"github.com/theartofdevel/notes_system/api_service/pkg/logging"
	"github.com/theartofdevel/notes_system/api_service/pkg/rest"
	"io/ioutil"
	"net/http"
	"strings"
	"time"
)

type Client struct {
	base rest.BaseClient
}

func NewClient(baseURL string, logger logging.Logger) *Client {
	return &Client{
		base: rest.BaseClient{
			BaseURL: baseURL,
			HTTPClient: &http.Client{
				Timeout: 10 * time.Second,
			},
			Logger: logger,
		},
	}
}

func (c *Client) GetCategories(userUuid string, ctx context.Context, filters []rest.FilterOptions) ([]byte, error) {
	var categories []byte

	c.base.Logger.Debug("add user_uuid to filter options")
	if filters == nil {
		filters = []rest.FilterOptions{
			{
				Field:  "user_uuid",
				Values: []string{userUuid},
			},
		}
	}

	c.base.Logger.Debug("build url with resource and filter")
	uri, err := rest.BuildURL(c.base.BaseURL, "/categories", filters)
	if err != nil {
		return categories, nil
	}
	c.base.Logger.Tracef("url: %s", uri)

	c.base.Logger.Debug("create new request")
	req, err := http.NewRequest("GET", uri, nil)
	if err != nil {
		return categories, errors.New(fmt.Sprintf("failed to create new request due to error: %v", err))
	}

	c.base.Logger.Debug("send request")
	req = req.WithContext(ctx)
	response, err := c.base.SendRequest(req)
	if err != nil {
		return categories, err
	}

	defer response.Body.Close()

	c.base.Logger.Debug("read body")
	categories, err = ioutil.ReadAll(response.Body)
	if err != nil {
		return nil, errors.New("failed to read body")
	}
	return categories, nil
}

// CreateCategory return new category uuid
func (c *Client) CreateCategory(categoryDTO CreateCategoryDTO, ctx context.Context, filters []rest.FilterOptions) (string, error) {
	var categoryUuid string

	c.base.Logger.Debug("build url with resource and filter")
	uri, err := rest.BuildURL(c.base.BaseURL, "/categories", filters)
	if err != nil {
		return categoryUuid, err
	}

	c.base.Logger.Debug("convert dto to map")
	structs.DefaultTagName = "json"
	data := structs.Map(categoryDTO)

	c.base.Logger.Debug("marshal map to bytes")
	dataBytes, err := json.Marshal(data)
	if err != nil {
		return categoryUuid, errors.New("failed to marshal dto")
	}

	c.base.Logger.Debug("create new request")
	req, err := http.NewRequest("POST", uri, bytes.NewBuffer(dataBytes))
	if err != nil {
		return categoryUuid, errors.New(fmt.Sprintf("failed to create new request due to error: %v", err))
	}

	c.base.Logger.Debug("send request")
	req = req.WithContext(ctx)
	response, err := c.base.SendRequest(req)
	if err != nil {
		return categoryUuid, err
	}

	c.base.Logger.Debug("parse location header")
	categoryURL, err := response.Location()
	if err != nil {
		return categoryUuid, errors.New("failed to get Location header")
	}
	c.base.Logger.Tracef("Location: %s", categoryURL.String())

	splitCategoryURL := strings.Split(categoryURL.String(), "/")
	categoryUuid = splitCategoryURL[len(splitCategoryURL)-1]
	return categoryUuid, nil
}

func (c *Client) UpdateCategory(categoryDTO UpdateCategoryDTO, ctx context.Context, filters []rest.FilterOptions) error {
	c.base.Logger.Debug("build url with resource and filter")
	uri, err := rest.BuildURL(c.base.BaseURL, fmt.Sprintf("/categories/%s", categoryDTO.Uuid), filters)
	if err != nil {
		return err
	}

	c.base.Logger.Debug("convert dto to map")
	structs.DefaultTagName = "json"
	data := structs.Map(categoryDTO)

	c.base.Logger.Debug("marshal map to bytes")
	dataBytes, err := json.Marshal(data)
	if err != nil {
		return errors.New("failed to marshal dto")
	}

	c.base.Logger.Debug("create new request")
	req, err := http.NewRequest("PATCH", uri, bytes.NewBuffer(dataBytes))
	if err != nil {
		return errors.New(fmt.Sprintf("failed to create new request due to error: %v", err))
	}

	c.base.Logger.Debug("send request")
	req = req.WithContext(ctx)
	_, err = c.base.SendRequest(req)
	if err != nil {
		return err
	}
	return nil
}

func (c *Client) DeleteCategory(categoryDTO DeleteCategoryDTO, ctx context.Context, filters []rest.FilterOptions) error {
	c.base.Logger.Debug("build url with resource and filter")
	uri, err := rest.BuildURL(c.base.BaseURL, fmt.Sprintf("/categories/%s", categoryDTO.Uuid), filters)
	if err != nil {
		return err
	}

	c.base.Logger.Debug("convert dto to map")
	structs.DefaultTagName = "json"
	data := structs.Map(categoryDTO)

	c.base.Logger.Debug("marshal map to bytes")
	dataBytes, err := json.Marshal(data)
	if err != nil {
		return errors.New("failed to marshal dto")
	}

	c.base.Logger.Debug("create new request")
	req, err := http.NewRequest("DELETE", uri, bytes.NewBuffer(dataBytes))
	if err != nil {
		return errors.New(fmt.Sprintf("failed to create new request due to error: %v", err))
	}

	c.base.Logger.Debug("send request")
	req = req.WithContext(ctx)
	_, err = c.base.SendRequest(req)
	if err != nil {
		return err
	}
	return nil
}
