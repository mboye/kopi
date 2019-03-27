package main

import (
	"encoding/base64"
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"os/signal"
	"path/filepath"
	"syscall"
	"unicode/utf8"

	"github.com/dustin/go-humanize"
	"github.com/mboye/kopi/model"
	log "github.com/sirupsen/logrus"
)

func main() {
	signals := make(chan os.Signal, 1)
	signal.Notify(signals, syscall.SIGINT, syscall.SIGTERM)

	recursive := flag.Bool("recursive", true, "Index path recursively")
	initIndex := flag.Bool("init", false, "Initial index. Mark all files as modified.")
	withProgress := flag.Bool("progress", true, "Print indexing progress.")
	flag.Parse()

	rootPath := flag.Arg(0)
	if rootPath == "" {
		log.Error("Path argument missing")
		printUsage()
		os.Exit(1)
	}

	log.Infof("Indexing path: %s", rootPath)
	log.Infof("Recursive: %t", *recursive)

	encoder := json.NewEncoder(os.Stdout)
	fileCount := int64(0)
	byteCount := int64(0)

	printProgress := func() {
		if fileCount > 0 && fileCount%1000 == 0 {
			log.WithFields(log.Fields{
				"files_found": fileCount,
				"bytes_found": humanize.Bytes(uint64(byteCount))}).Info("Progress")
		}
	}

	walkFn := func(path string, info os.FileInfo, err error) error {
		convertedPath := false
		if !utf8.ValidString(path) {
			log.WithField("path", path).Warn("path is invalid utf8 string; converting path to base64")
			path = base64.RawStdEncoding.EncodeToString([]byte(path))
			convertedPath = true
		}

		if os.IsPermission(err) {
			log.WithField("path", path).Warn("Permission denied")
		} else if err != nil {
			log.Errorf("Failed to walk path: %s", path)
			return err
		}

		select {
		case sig := <-signals:
			log.Fatalf("Received signal: %s", sig.String())
		default:
		}

		log.Debugf("Walking path: %s", path)

		if info.IsDir() && path != rootPath && !*recursive {
			return filepath.SkipDir
		}

		if !info.IsDir() && !info.Mode().IsRegular() {
			log.WithField("path", path).Debug("Ignoring non-regular file")
			return nil
		}

		size := int64(0)
		if !info.IsDir() {
			size = info.Size()
		}

		file := &model.File{
			Path:          path,
			ConvertedPath: convertedPath,
			Size:          size,
			Mode:          info.Mode(),
			ModifiedTime:  info.ModTime().UTC()}

		if *initIndex {
			file.Modified = true
		}

		if err := encoder.Encode(file); err != nil {
			log.Errorf("Failed to encode file: %s", path)
			return err
		}
		fileCount++
		byteCount += size

		if *withProgress {
			printProgress()
		}

		return nil
	}

	if err := filepath.Walk(rootPath, walkFn); err != nil {
		log.Errorf("Indexing failed: %s", err.Error())
		os.Exit(1)
	}

	log.Infof("Number of files indexed: %d", fileCount)
	log.Infof("Number of bytes indexed: %s", humanize.Bytes(uint64(byteCount)))

}

func printUsage() {
	commandName := filepath.Base(os.Args[0])
	fmt.Printf("Usage: %s <path>\n", commandName)
}
