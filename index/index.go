package index

import "github.com/mboye/kopi/model"

type WalkFunc func(path string, file *model.File) error

type Index interface {
	Find(path string) *model.File
	Add(file *model.File) error
	Print()
	Walk(walkFn WalkFunc) error
	Size() int
}
