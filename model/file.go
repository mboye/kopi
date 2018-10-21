package model

import (
	"os"
	"time"
)

type File struct {
	Path         string      `json:"path"`
	Size         int64       `json:"size"`
	ModifiedTime time.Time   `json:"modifiedTime"`
	Mode         os.FileMode `json:"mode"`
}
