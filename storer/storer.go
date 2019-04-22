package storer

import (
	"errors"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/mboye/kopi/input"
	"github.com/mboye/kopi/model"
	"github.com/mboye/kopi/outputhandler"
	"github.com/mboye/kopi/security"
	"github.com/mboye/kopi/stage"
	log "github.com/sirupsen/logrus"
)

type storer struct {
	outputDir        string
	maxBlockSize     int64
	encrypt          bool
	progressInterval uint
}

var _ stage.Stage = (*storer)(nil)

func New(outputDir string, maxBlockSize int64, encrypt bool, progressInterval uint) (stage.Stage, error) {
	if outputDir == "" {
		return nil, errors.New("cannot store to empty output dir")
	}

	if maxBlockSize < 0 {
		return nil, errors.New("max block size must be > 0")
	}

	if progressInterval < 1 {
		return nil, errors.New("progress interval must be >= 1")
	}

	return &storer{
		outputDir,
		maxBlockSize, encrypt, progressInterval}, nil
}

func (s *storer) Execute() error {
	securityContext, err := security.NewContext(s.outputDir, s.encrypt)
	if err != nil {
		log.WithField("error", err).Fatal("failed to create security context")
	}

	filterAndStoreFile := func(file *model.File) error {
		if file.Mode.IsDir() {
			outputhandler.Stdout.Handle(file)
			return nil
		}

		if !file.Modified {
			log.WithField("path", file.Path).Debug("Skipping unmodified file")
			outputhandler.Stdout.Handle(file)
			return nil
		}

		err := storeFile(file, s.outputDir, securityContext, s.maxBlockSize)
		if os.IsNotExist(err) {
			log.WithField("path", file.Path).Warn("File not found")
			return nil
		} else if os.IsPermission(err) {
			log.WithField("path", file.Path).Warn("Permission denied")
			return nil
		}
		return err
	}

	log.WithField("destination", s.outputDir).Info("Beginning to store files")
	return input.ProcessFilesWithProgress(filterAndStoreFile, s.progressInterval)
}

func storeFile(file *model.File, outputDir string, securityContext *security.Context, maxBlockSize int64) error {
	if err := refreshFileMetadata(file); err != nil {
		return err
	}

	logger := log.WithField("path", file.Path)

	inputFile, err := os.Open(file.Path)
	if err != nil {
		return err
	}
	defer inputFile.Close()
	fileReader := io.LimitReader(inputFile, file.Size)

	logger.WithField("max_offset", file.Size).Debug("Read file")

	fileOffset := int64(0)
	for fileOffset < file.Size {
		blockReader := io.LimitReader(fileReader, maxBlockSize)
		blockOffset := fileOffset

		hasher, err := securityContext.NewHasher()
		if err != nil {
			log.WithField("error", err).Fatal("failed to get hasher")
		}

		blockData, err := ioutil.ReadAll(blockReader)
		bytesRead := len(blockData)
		logger.WithField("bytes_read", bytesRead).Debug("Read file")

		fileOffset += int64(bytesRead)

		if int64(bytesRead) != maxBlockSize && fileOffset != file.Size {
			return errors.New("Incomplete buffer read")
		}

		blockSize := bytesRead
		hasher.Write(blockData)
		hash := fmt.Sprintf("%x", hasher.Sum(nil))
		block := model.Block{Hash: hash, Offset: blockOffset, Size: int64(blockSize)}

		outputPath := fmt.Sprintf("%s/%s/%s.block", outputDir, hash[:2], hash)
		_, err = os.Stat(outputPath)
		if err == nil {
			logger.WithFields(log.Fields{"hash": hash, "offset": blockOffset, "size": blockSize}).Debug("Reusing existing block")
			file.AddBlock(block)
			continue
		}

		encodedBlockData, err := securityContext.Encode(blockData)
		if err != nil {
			return err
		}

		blockDirPath := fmt.Sprintf("%s/%s", outputDir, hash[:2])
		err = os.MkdirAll(blockDirPath, 0755)
		if err != nil {
			return err
		}

		outputFile, err := os.Create(outputPath)
		if err != nil {
			return err
		}
		logger.WithField("output_path", outputPath).Debug("Created output file")

		bytesWritten, err := outputFile.Write(encodedBlockData)
		if err != nil {
			outputFile.Close()
			return err
		}

		if bytesWritten != len(encodedBlockData) {
			return errors.New("Incomplete block write")
		}

		outputFile.Close()
		file.AddBlock(block)
		logger.WithFields(log.Fields{"hash": hash, "offset": blockOffset, "size": blockSize}).Debug("Block created")

		logger.WithField("fileOffset", fileOffset).Debug("File offset")
	}

	logger.Debug("File read completed")
	outputhandler.Stdout.Handle(file)
	return nil
}

func refreshFileMetadata(file *model.File) error {
	if fileInfo, err := os.Stat(file.Path); err != nil {
		log.WithError(err).Error("Failed to get file metadata")
		return err
	} else {
		if file.Size != fileInfo.Size() {
			log.WithFields(log.Fields{"path": file.Path, "expected_size": file.Size, "actual_size": fileInfo.Size()}).Warn("File size has changed")
		}
		file.Size = fileInfo.Size()

		if file.ModifiedTime.UTC() != fileInfo.ModTime().UTC() {
			log.WithFields(log.Fields{"path": file.Path, "expected_modtime": file.ModifiedTime, "actual_modtime": fileInfo.ModTime()}).Warn("File modification time has changed")
		}
		file.ModifiedTime = fileInfo.ModTime().UTC()

		if file.Mode != fileInfo.Mode() {
			log.WithFields(log.Fields{"path": file.Path, "expected_mode": file.Mode, "actual_mode": fileInfo.Mode()}).Warn("File mode has changed")
		}
		file.Mode = fileInfo.Mode()
	}

	return nil
}

func printUsage() {
	commandName := filepath.Base(os.Args[0])
	fmt.Fprintf(os.Stderr, "Usage: %s [OPTIONS] <destination dir>\n\n", commandName)
	fmt.Fprintf(os.Stderr, "Pass index lines on STDIN.\n")
	fmt.Fprintln(os.Stderr, "\nOptions:")
	flag.PrintDefaults()
}
