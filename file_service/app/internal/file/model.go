package file

import (
	"fmt"
	"io"
	"io/ioutil"
)

type File struct {
	Name  string `json:"name"`
	Size  int64  `json:"size"`
	Bytes []byte `json:"file"`
}

type CreateFileDTO struct {
	Name     string `json:"name"`
	Size     int64  `json:"size"`
	Reader   io.Reader
}

func NewFile(dto CreateFileDTO) (*File, error) {
	bytes, err := ioutil.ReadAll(dto.Reader)
	if err != nil {
		return nil, fmt.Errorf("failed to create file model. err: %w", err)
	}
	return &File{
		Name:  dto.Name,
		Size:  dto.Size,
		Bytes: bytes,
	}, nil
}
