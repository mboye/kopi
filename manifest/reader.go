package manifest

import (
	"bytes"
	"compress/gzip"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"os"

	_ "github.com/mboye/kopi/loglevel"
	"github.com/mboye/kopi/model"
	"github.com/mboye/kopi/outputhandler"
	"github.com/mboye/kopi/security"
	"github.com/mboye/kopi/stage"
	log "github.com/sirupsen/logrus"
)

type reader struct {
	inputDir string
	decrypt  bool
	id       string
}

var _ stage.Stage = (*reader)(nil)

func NewReader(inputDir string, decrypt bool, id string) (stage.Stage, error) {
	if inputDir == "" {
		return nil, errors.New("input dir cannot be empty")
	}

	if id == "" {
		return nil, errors.New("cannot read manifest with empty ID")
	}

	return &reader{inputDir, decrypt, id}, nil
}

func (r *reader) Execute() error {
	manifestPath := fmt.Sprintf("%s/manifests/%s", r.inputDir, r.id)
	manifestFile, err := os.Open(manifestPath)
	if err != nil {
		return fmt.Errorf("failed to open manifest: %s", err.Error())
	}
	defer manifestFile.Close()

	securityContext, err := security.NewContext(r.inputDir, r.decrypt)
	if err != nil {
		return fmt.Errorf("failed to create security context: %s", err.Error())
	}

	encodedData, err := ioutil.ReadAll(manifestFile)
	if err != nil {
		return fmt.Errorf("failed to decompress manifest: %s", err.Error())
	}

	compressedData, err := securityContext.Decode(encodedData)
	if err != nil {
		return fmt.Errorf("failed to decode manifest: %s", err.Error())
	}

	decompressor, err := gzip.NewReader(bytes.NewBuffer(compressedData))
	if err != nil {
		return fmt.Errorf("failed to create decompressor: %s", err.Error())
	}
	defer decompressor.Close()

	decoder := json.NewDecoder(decompressor)
	decoder.DisallowUnknownFields()

	header := manifestHeader{}
	if err := decoder.Decode(&header); err != nil {
		return fmt.Errorf("failed to decode manifest header: %s", err.Error())
	}
	log.WithFields(log.Fields{
		"date":        header.Date,
		"id":          header.ID,
		"description": header.Description,
	}).Info("read manifest header")

	for decoder.More() {
		file := model.File{}
		if err := decoder.Decode(&file); err != nil {
			return fmt.Errorf("failed to decode file: %s", err.Error())
		}

		outputhandler.Stdout.Handle(file)
	}

	return nil
}
