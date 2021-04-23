package file

import (
	"fmt"
	"github.com/google/uuid"
	"golang.org/x/text/transform"
	"golang.org/x/text/unicode/norm"
	"io"
	"io/ioutil"
	"strings"
	"unicode"
)

type File struct {
	ID 	  string `json:"id"`
	Name  string `json:"name"`
	Size  int64  `json:"size"`
	Bytes []byte `json:"file"`
}

type CreateFileDTO struct {
	Name     string `json:"name"`
	Size     int64  `json:"size"`
	Reader   io.Reader
}

func isMn(r rune) bool {
	return unicode.Is(unicode.Mn, r) // Mn: nonspacing marks
}

func (d CreateFileDTO) NormalizeName() {
	d.Name = strings.ReplaceAll(d.Name, " ", "_")
	t := transform.Chain(norm.NFD, transform.RemoveFunc(isMn), norm.NFC)
	d.Name, _, _ = transform.String(t, d.Name)
}

func NewFile(dto CreateFileDTO) (*File, error) {
	bytes, err := ioutil.ReadAll(dto.Reader)
	if err != nil {
		return nil, fmt.Errorf("failed to create file model. err: %w", err)
	}
	id, err := uuid.NewUUID()
	if err != nil {
		return nil, fmt.Errorf("failed to generate file id. err: %w", err)
	}

	return &File{
		ID:    id.String(),
		Name:  dto.Name,
		Size:  dto.Size,
		Bytes: bytes,
	}, nil
}
