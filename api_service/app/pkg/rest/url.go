package rest

import (
	"errors"
	"fmt"
	"net/url"
	"path"
	"strings"
)

type FilterOptions struct {
	Field    string
	Operator string
	Values   []string
}

// filtering options like in:1,3,4 or neq:4 or eq:1 or =123
func (fo *FilterOptions) ToStringWF() string {
	return fmt.Sprintf("%s%s", fo.Operator, strings.Join(fo.Values, ","))
}

func BuildURL(baseURL string, resource string, filters []FilterOptions) (string, error) {
	var resultURL string
	parsedURL, err := url.ParseRequestURI(baseURL)
	if err != nil {
		return resultURL, errors.New(fmt.Sprintf("failed to parse base URL. error: %v", err))
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

type APIErrorResponse struct {
	Error            string `json:"error"`
	ErrorCode        string `json:"error_code"`
	DeveloperMessage string `json:"developer_message"`
}

func (aep *APIErrorResponse) ToString() string {
	return fmt.Sprintf("Error Code: %s, Error: %s, Developer Message: %s", aep.ErrorCode, aep.Error, aep.DeveloperMessage)
}
