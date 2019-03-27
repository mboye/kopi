package model

import (
	"os"
	"path/filepath"
	"time"
)

type File struct {
	Path          string      `json:"path"`
	ConvertedPath bool        `json:"convertedPath,omitEmpty"`
	Size          int64       `json:"size"`
	ModifiedTime  time.Time   `json:"modifiedTime"`
	Mode          os.FileMode `json:"mode"`
	Modified      bool        `json:"modified,omitempty"`
	Blocks        []Block     `json:"blocks,omitempty"`
}

type Block struct {
	Hash   string `json:"hash"`
	Offset int64  `json:"offset"`
	Size   int64  `json:"size"`
}

func NewFile() *File {
	return &File{Blocks: []Block{}}
}

func (f *File) Name() string {
	return filepath.Base(f.Path)
}

func (f *File) AddBlock(block Block) {
	f.Blocks = append(f.Blocks, block)
}

func FilesEqual(a, b *File) bool {
	return a.Path == b.Path && a.Size == b.Size && a.ModifiedTime.UTC() == b.ModifiedTime.UTC() && a.Mode == b.Mode
}
