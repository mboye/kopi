package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"path/filepath"

	"github.com/mboye/kopi/index"
	"github.com/mboye/kopi/model"
	"github.com/mboye/kopi/util"
	log "github.com/sirupsen/logrus"
)

func main() {
	util.SetLogLevel()
	flag.Parse()

	if flag.NArg() != 2 {
		log.Error("Path argument missing")
		printUsage()
		os.Exit(1)
	}

	pathA := flag.Arg(0)
	pathB := flag.Arg(1)

	log.WithFields(log.Fields{"index_a": pathA, "index_b": pathB}).Info("Diffing indices")

	var indexA, indexB index.Index
	var err error
	if indexA, err = loadIndex(pathA); err != nil {
		log.Fatal(err)
	}
	log.WithField("size", indexA.Size()).Info("Loaded index A")

	if indexB, err = loadIndex(pathB); err != nil {
		log.Fatal(err)
	}
	log.WithField("size", indexB.Size()).Info("Loaded index B")

	markIndexChanges(indexA, indexB)
	indexB.Print()
}

func loadIndex(path string) (index.Index, error) {
	log.Debugf("Loading index: %s", path)
	inputFile, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("failed to open index: %s", err.Error())
	}
	decoder := json.NewDecoder(inputFile)

	index := index.New()
	for decoder.More() {
		file := &model.File{}
		if err = decoder.Decode(file); err != nil {
			return nil, fmt.Errorf("failed to decode file: %s", err.Error())
		}

		file.Modified = false
		if err := index.Add(file); err != nil {
			return nil, err
		}
	}
	return index, nil
}

func markIndexChanges(indexA, indexB index.Index) {
	modifiedCount := int64(0)
	diffWalker := func(pathB string, fileB *model.File) error {
		if fileA := indexA.Find(pathB); fileA != nil {
			if !model.FilesEqual(fileA, fileB) {
				fileB.Modified = true
				modifiedCount++
			} else {
				// Preserve blocks of unmodified file
				fileB.Blocks = fileA.Blocks
			}
		} else {
			fileB.Modified = true
			modifiedCount++
		}

		return nil
	}

	indexB.Walk(diffWalker)

	log.WithField("modified_files", modifiedCount).Info("Diffing completed")
}

func printUsage() {
	commandName := filepath.Base(os.Args[0])
	fmt.Printf("Usage: %s <index a> <index b>\n", commandName)
}
