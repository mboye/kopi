package input

import (
	"encoding/json"
	"fmt"
	"math"
	"os"
	"time"

	"github.com/dustin/go-humanize"
	"github.com/mboye/kopi/model"
	log "github.com/sirupsen/logrus"
)

type FileHandlerFunc func(file *model.File) error

func ProcessFiles(handler FileHandlerFunc) error {
	decoder := json.NewDecoder(os.Stdin)

	for decoder.More() {
		file := &model.File{}
		if err := decoder.Decode(file); err != nil {
			return fmt.Errorf("Failed to decode file: %s", err.Error())
		}

		if err := handler(file); err != nil {
			return err
		}

	}
	return nil
}

func ProcessFilesWithProgress(handler FileHandlerFunc) error {
	files := []*model.File{}
	var maxFiles, maxBytes int64
	var filesProcessed, bytesProcessed, errorCount int64
	startTime := time.Now()

	addToSummary := func(f *model.File) error {
		files = append(files, f)
		maxFiles++
		maxBytes += f.Size
		return nil
	}

	if err := ProcessFiles(addToSummary); err != nil {
		return err
	}

	printProgress := func() {
		if bytesProcessed > maxBytes {
			bytesProcessed = maxFiles
		}

		fileProgress := 100.0 * float32(filesProcessed) / float32(maxFiles)
		byteProgress := 100.0 * float32(bytesProcessed) / float32(maxBytes)
		elapsedTime := time.Now().Sub(startTime).Round(time.Second)

		fileRate := float64(filesProcessed) / elapsedTime.Seconds()
		byteRate := float64(bytesProcessed) / elapsedTime.Seconds()
		fileRemainingTime := int64(math.Round(float64(maxFiles-filesProcessed) / fileRate))
		byteRemainingTime := int64(math.Round(float64(maxBytes-bytesProcessed) / byteRate))

		var remainingTime time.Duration
		if fileRemainingTime > byteRemainingTime {
			remainingTime = time.Duration(fileRemainingTime * 1e9)
		} else {
			remainingTime = time.Duration(byteRemainingTime * 1e9)
		}

		log.WithFields(log.Fields{
			"file_progress": fmt.Sprintf("%d / %d = %.2f%%", filesProcessed, maxFiles, fileProgress),
			"byte_progress": fmt.Sprintf("%s / %s = %.2f%%",
				humanize.Bytes(uint64(bytesProcessed)), humanize.Bytes(uint64(maxBytes)), byteProgress),
			"elapsed_time":   elapsedTime.String(),
			"remaining_time": remainingTime.String(),
			"errors":         errorCount}).Info("progress")
	}

	for _, file := range files {
		if err := handler(file); err != nil {
			return err
		}

		filesProcessed++
		bytesProcessed += file.Size
		printProgress()
	}

	return nil
}