package rest

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/theartofdevel/notes_system/api_service/pkg/logging"
	"net/http"
	"net/url"
	"path"
)

type BaseClient struct {
	BaseURL    string
	HTTPClient *http.Client
	Logger     logging.Logger
}

func (c *BaseClient) SendRequest(req *http.Request) (*APIResponse, error) {
	if c.HTTPClient == nil {
		return nil, errors.New("no http client")
	}

	req.Header.Set("Accept", "application/json; charset=utf-8")
	req.Header.Set("Content-Type", "application/json; charset=utf-8")

	response, err := c.HTTPClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to send request. error: %w", err)
	}

	apiResponse := APIResponse{
		IsOk: true,
		response: response,
	}
	if response.StatusCode < http.StatusOK || response.StatusCode >= http.StatusBadRequest {
		apiResponse.IsOk = false
		// if an error, read body so close it
		defer response.Body.Close()

		var apiErr APIError
		if err = json.NewDecoder(response.Body).Decode(&apiErr); err == nil {
			apiResponse.Error = apiErr
		}
	}

	return &apiResponse, nil
}

func (c *BaseClient) BuildURL(resource string, filters []FilterOptions) (string, error) {
	var resultURL string
	parsedURL, err := url.ParseRequestURI(c.BaseURL)
	if err != nil {
		return resultURL, fmt.Errorf("failed to parse base URL. error: %w", err)
	}
	parsedURL.Path = path.Join(parsedURL.Path, resource)

	if len(filters) > 0 {
		q := parsedURL.Query()
		for _, fo := range filters {
			q.Set(fo.Field, fo.ToStringWF())
		}
		parsedURL.RawQuery = q.Encode()
	}

	return parsedURL.String(), nil
}

func (c *BaseClient) Close() error {
	c.HTTPClient = nil
	return nil
}
