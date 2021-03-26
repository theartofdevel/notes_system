package rest

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/theartofdevel/notes_system/api_service/pkg/logging"
	"net/http"
	"sync"
)

type BaseClient struct {
	BaseURL    string
	HTTPClient *http.Client
	mu         sync.Mutex
	Logger     logging.Logger
}

func (c *BaseClient) SendRequest(req *http.Request) (*http.Response, error) {
	c.mu.Lock()
	defer c.mu.Unlock()

	if c.HTTPClient == nil {
		return nil, errors.New("no http client")
	}

	req.Header.Set("Accept", "application/json; charset=utf-8")
	req.Header.Set("Content-Type", "application/json; charset=utf-8")

	response, err := c.HTTPClient.Do(req)
	if err != nil {
		return nil, errors.New(fmt.Sprintf("failed to send request. error: %s", err))
	}

	if response.StatusCode < http.StatusOK || response.StatusCode >= http.StatusBadRequest {
		// if an error, read body so close it
		defer response.Body.Close()

		var errResponse APIErrorResponse
		if err = json.NewDecoder(response.Body).Decode(&errResponse); err == nil {
			return nil, errors.New(errResponse.ToString())
		}

		return nil, errors.New(fmt.Sprintf("got an error response, can't decode to error response, StatusCode: %d, Body: %v", response.StatusCode, response.Body))
	}

	return response, nil
}


func (c *BaseClient) Close() error {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.HTTPClient = nil
	return nil
}