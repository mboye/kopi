package manifest

import (
	"bytes"
	"compress/gzip"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"time"

	"github.com/dustin/go-humanize"
	"github.com/mboye/kopi/input"
	_ "github.com/mboye/kopi/loglevel"
	"github.com/mboye/kopi/model"
	"github.com/mboye/kopi/security"
	"github.com/mboye/kopi/stage"
	log "github.com/sirupsen/logrus"
)

type writer struct {
	outputDir   string
	encrypt     bool
	description string
}

var _ stage.Stage = (*writer)(nil)

func NewWriter(outputDir string, encrypt bool, description string) (stage.Stage, error) {

	if outputDir == "" {
		return nil, errors.New("output directory cannot be empty")
	}

	return &writer{outputDir, encrypt, description}, nil
}
func (w *writer) Execute() error {
	now := time.Now().UTC()
	manifestFilename := fmt.Sprintf("%d/%02d/%02d/%d.manifest",
		now.Year(), now.Month(), now.Day(),
		now.Unix())

	header := manifestHeader{
		ID:          manifestFilename,
		Date:        now,
		Description: w.description}

	securityContext, err := security.NewContext(w.outputDir, w.encrypt)
	if err != nil {
		log.WithField("error", err).Fatal("failed to create security context")
	}

	compressedManifest := bytes.NewBuffer(nil)
	compressor := gzip.NewWriter(compressedManifest)

	addToManifest := func(item interface{}) error {
		data, err := json.Marshal(item)
		if err != nil {
			return err
		}

		_, err = io.CopyN(compressor, bytes.NewBuffer(data), int64(len(data)))
		if err != nil {
			return err
		}

		_, err = compressor.Write([]byte{'\n'})
		return err
	}

	var fileCount, byteCount int64
	addFileToManifest := func(file *model.File) error {
		fileCount++
		byteCount += file.Size
		return addToManifest(file)
	}

	addToManifest(header)

	err = input.ProcessFilesWithProgress(addFileToManifest, 1)
	if err != nil {
		log.WithError(err).Error("failed to create compressed manifest")
		return err
	}
	compressor.Close()

	log.WithField("size", humanize.Bytes(uint64(compressedManifest.Len()))).Info("compressed manifest created")

	encodedManifest, err := securityContext.Encode(compressedManifest.Bytes())
	if err != nil {
		return err
	}
	log.WithField("size", len(encodedManifest)).Debug("encoded compressed manifest")

	outputPath := fmt.Sprintf("%s/manifests/%s", w.outputDir, manifestFilename)
	if err := os.MkdirAll(filepath.Dir(outputPath), 0755); err != nil {
		log.Error("failed to create manifest directory")
		return err
	}

	outputFile, err := os.OpenFile(outputPath, os.O_WRONLY|os.O_CREATE|os.O_EXCL, 0644)
	if err != nil {
		return err
	}
	defer outputFile.Close()
	log.WithField("path", outputFile.Name()).Debug("opened manifest file")

	bytesWritten, err := io.CopyN(outputFile, bytes.NewBuffer(encodedManifest), int64(len(encodedManifest)))
	if err != nil {
		return fmt.Errorf("failed to save manifest: %s", err.Error())
	}

	log.WithField("bytes_written", bytesWritten).Debug("wrote manifest")

	log.WithFields(log.Fields{
		"id":    header.ID,
		"files": fileCount,
		"bytes": humanize.Bytes(uint64(byteCount))}).Info("created manifest")

	return nil
}
