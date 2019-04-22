package scanner

import (
	"errors"
	"fmt"
	"os"
	"os/signal"
	"path/filepath"
	"syscall"
	"unicode/utf8"

	"github.com/dustin/go-humanize"
	"github.com/mboye/kopi/model"
	oh "github.com/mboye/kopi/outputhandler"
	"github.com/mboye/kopi/stage"
	log "github.com/sirupsen/logrus"
)

type scanner struct {
	rootPath     string
	recursive    bool
	initial      bool
	withProgress bool
}

var _ stage.Stage = (*scanner)(nil)

func New(rootPath string, recursive bool, initial bool, withProgress bool) (stage.Stage, error) {
	if rootPath == "" {
		return nil, errors.New("cannot scan empty path")
	}

	return &scanner{
		rootPath:  rootPath,
		recursive: recursive,
		initial:   initial, withProgress: withProgress}, nil
}

func (s *scanner) Execute() error {
	signals := make(chan os.Signal, 1)
	signal.Notify(signals, syscall.SIGINT, syscall.SIGTERM)
	defer signal.Reset(syscall.SIGINT, syscall.SIGTERM)

	log.Infof("Indexing path: %s", s.rootPath)
	log.Infof("Recursive: %t", s.recursive)

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
		if !utf8.ValidString(path) {
			log.WithField("path", path).Warn("ignoring file with non-utf8 path")
			return nil
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

		if info.IsDir() && path != s.rootPath && !s.recursive {
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
			Path:         path,
			Size:         size,
			Mode:         info.Mode(),
			ModifiedTime: info.ModTime().UTC()}

		if s.initial {
			file.Modified = true
		}

		if err := oh.Stdout.Handle(file); err != nil {
			return fmt.Errorf("output handler failed: %s", err.Error())
		}
		fileCount++
		byteCount += size

		if s.withProgress {
			printProgress()
		}

		return nil
	}

	if err := filepath.Walk(s.rootPath, walkFn); err != nil {
		return fmt.Errorf("indexing failed: %s", err.Error())
	}

	log.Infof("Number of files indexed: %d", fileCount)
	log.Infof("Number of bytes indexed: %s", humanize.Bytes(uint64(byteCount)))
	return nil
}
