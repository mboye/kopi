package input

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"math"
	"os"
	"os/signal"
	"sync/atomic"
	"syscall"
	"time"

	"github.com/dustin/go-humanize"
	"github.com/mboye/kopi/model"
	log "github.com/sirupsen/logrus"
)

type FileHandlerFunc func(file *model.File) error

func ProcessFiles(handler FileHandlerFunc) error {
	signals := make(chan os.Signal, 1)
	signal.Notify(signals, syscall.SIGINT, syscall.SIGTERM)

	decoder := json.NewDecoder(os.Stdin)

	for decoder.More() {
		select {
		case sig := <-signals:
			log.Fatalf("Received signal: %s", sig.String())
		default:
		}

		file := &model.File{}
		if err := decoder.Decode(file); err != nil {
			return fmt.Errorf("Failed to decode file: %s", err.Error())
		}

		if file.ConvertedPath {
			originalPath, err := base64.RawStdEncoding.DecodeString(file.Path)
			if err != nil {
				return fmt.Errorf("Failed to decode file path: %s", err.Error())
			}
			file.Path = string(originalPath)
		}

		if err := handler(file); err != nil {
			return err
		}

	}
	return nil
}

func ProcessFilesWithProgress(handler FileHandlerFunc, interval uint) error {
	files := []*model.File{}
	var maxFiles, maxBytes int64
	var filesProcessed, bytesProcessed int64
	startTime := time.Now()

	signals := make(chan os.Signal, 1)
	signal.Notify(signals, syscall.SIGINT, syscall.SIGTERM)

	addToSummary := func(f *model.File) error {
		files = append(files, f)
		maxFiles++
		maxBytes += f.Size
		return nil
	}

	if err := ProcessFiles(addToSummary); err != nil {
		return err
	}

	progressPrinter := func(stop chan struct{}) {
		ticker := time.NewTicker(time.Duration(interval * 1e9))
		for {
			select {
			case <-ticker.C:
				filesProcessedSnapshot := atomic.LoadInt64(&filesProcessed)
				bytesProcessedSnapshot := atomic.LoadInt64(&bytesProcessed)
				printProgress(maxFiles, maxBytes, filesProcessedSnapshot, bytesProcessedSnapshot, startTime)
			case <-stop:
				return
			}
		}

	}

	if interval > 0 {
		printProgress(maxFiles, maxBytes, filesProcessed, bytesProcessed, startTime)

		stop := make(chan struct{}, 1)
		go progressPrinter(stop)
		defer func() {
			close(stop)
		}()
	}

	for _, file := range files {
		select {
		case sig := <-signals:
			log.Fatalf("Received signal: %s", sig.String())
		default:
		}

		if err := handler(file); err != nil {
			return err
		}

		atomic.AddInt64(&filesProcessed, 1)
		atomic.AddInt64(&bytesProcessed, file.Size)
	}
	printProgress(maxFiles, maxBytes, filesProcessed, bytesProcessed, startTime)

	return nil
}

func printProgress(maxFiles, maxBytes, filesProcessed, bytesProcessed int64, startTime time.Time) {
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
	}).Info("Progress")

}
