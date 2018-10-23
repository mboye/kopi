package index

import (
	"encoding/json"
	"errors"
	"os"
	"sort"

	"github.com/mboye/kopi/model"
)

type simpleIndex struct {
	files map[string]*model.File
}

func New() Index {
	return &simpleIndex{files: make(map[string]*model.File, 0)}
}

func (t *simpleIndex) Add(file *model.File) error {
	if t.Find(file.Path) != nil {
		return errors.New("file path already in Index")
	}

	t.files[file.Path] = file
	return nil
}

func (t *simpleIndex) Find(path string) *model.File {
	if file, found := t.files[path]; found {
		return file
	}
	return nil
}

func (t *simpleIndex) Print() {
	sortedFiles := make([]*model.File, len(t.files))
	i := 0
	for _, file := range t.files {
		sortedFiles[i] = file
		i++
	}

	filePathLess := func(a, b int) bool {
		pathA := sortedFiles[a].Path
		pathB := sortedFiles[b].Path
		return pathA < pathB
	}

	sort.Slice(sortedFiles, filePathLess)

	encoder := json.NewEncoder(os.Stdout)
	for _, file := range sortedFiles {
		encoder.Encode(file)
	}
}

func (t *simpleIndex) Walk(walkFn WalkFunc) error {
	for path, file := range t.files {
		if err := walkFn(path, file); err != nil {
			return err
		}
	}
	return nil
}

func (t *simpleIndex) Size() int {
	return len(t.files)
}
