package model

import (
	"os"
	"path/filepath"
	"time"
)

type File struct {
	Path         string      `json:"path"`
	Size         int64       `json:"size"`
	ModifiedTime time.Time   `json:"modifiedTime"`
	Mode         os.FileMode `json:"mode"`
	Modified     bool        `json:"modified,omitempty"`
}

func NewFile() *File {
	return &File{}
}

func (f *File) Name() string {
	return filepath.Base(f.Path)
}
