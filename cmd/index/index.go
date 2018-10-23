package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"path/filepath"

	"github.com/mboye/kopi/model"
	log "github.com/sirupsen/logrus"
)

func main() {
	recursive := true
	flag.BoolVar(&recursive, "recursive", true, "Index path recursively")
	flag.Parse()

	rootPath := flag.Arg(0)
	if rootPath == "" {
		log.Error("Path argument missing")
		printUsage()
		os.Exit(1)
	}

	log.Infof("Indexing path: %s", rootPath)
	log.Infof("Recursive: %t", recursive)

	encoder := json.NewEncoder(os.Stdout)
	fileCount := 0

	walkFn := func(path string, info os.FileInfo, err error) error {
		if err != nil {
			log.Errorf("Failed to walk path: %s", path)
			return err
		}

		log.Debugf("Walking path: %s", path)

		if info.IsDir() && path != rootPath && !recursive {
			return filepath.SkipDir
		}

		size := int64(0)
		if !info.IsDir() {
			size = info.Size()
		}

		file := &model.File{
			Path:         path,
			Size:         size,
			Mode:         info.Mode(),
			ModifiedTime: info.ModTime().UTC()}

		if err := encoder.Encode(file); err != nil {
			log.Errorf("Failed to encode file: %s", path)
			return err
		}
		fileCount++
		return nil
	}

	if err := filepath.Walk(rootPath, walkFn); err != nil {
		log.Errorf("Indexing failed: %s", err.Error())
		os.Exit(1)
	}

	log.Infof("Number of files indexed: %d", fileCount)

}

func printUsage() {
	commandName := filepath.Base(os.Args[0])
	fmt.Printf("Usage: %s <path>\n", commandName)
}
