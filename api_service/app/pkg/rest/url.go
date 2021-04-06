package rest

import (
	"fmt"
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
