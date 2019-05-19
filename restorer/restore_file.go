package restorer

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/mboye/kopi/model"
	"github.com/mboye/kopi/security"
	log "github.com/sirupsen/logrus"
)

func restoreFile(file *model.File, inputDir, outputDir string, securityContext *security.Context, dryRun bool) error {
	log.WithFields(log.Fields{
		"path": file.Path,
		"mode": file.Mode}).Debug("restoring file")

	if file.Size > 0 && (file.Blocks == nil || len(file.Blocks) == 0) {
		return errors.New("cannot restore non-empty file without blocks")
	}

	outputPath := fmt.Sprintf("%s/%s", outputDir, file.Path)

	var outputFile *os.File
	if !dryRun {
		parentDir := filepath.Dir(outputPath)
		if err := os.MkdirAll(parentDir, 0755); err != nil {
			log.WithError(err).Error("failed to create parent directory of file")
			return err
		}

		var err error
		outputFile, err = os.OpenFile(outputPath, os.O_WRONLY|os.O_CREATE|os.O_EXCL, file.Mode)
		if err != nil {
			return err
		}
		defer outputFile.Close()

		if file.Size == 0 {
			return nil
		}
	}

	progress := 0
	for _, block := range file.Blocks {
		progress++
		log.WithFields(
			log.Fields{
				"progress":     progress,
				"max_progress": len(file.Blocks)}).Debug("restoring block")

		restoreBlock := func() error {
			blockPath := fmt.Sprintf("%s/%s/%s.block", inputDir, block.Hash[:2], block.Hash)
			blockFile, err := os.Open(blockPath)
			if err != nil {
				log.WithError(err).Error("failed to open block file")
				return err
			}
			defer blockFile.Close()
			log.WithField("path", blockPath).Debug("opened block")

			encodedBlockData, err := ioutil.ReadFile(blockPath)
			if err != nil {
				log.WithError(err).Error("failed to read block data")
				return err
			}
			log.WithField("size", len(encodedBlockData)).Debug("read block data")

			blockData, err := securityContext.Decode(encodedBlockData)
			actualBlockSize := len(blockData)
			log.WithField("size", actualBlockSize).Debug("decoded block data")

			if actualBlockSize < int(block.Size) {
				log.WithFields(
					log.Fields{
						"actual_size":       actualBlockSize,
						"min_expected_size": block.Size,
						"block_path":        blockPath}).Error("corrupt block detected")

				return fmt.Errorf("corrupt block: %s", blockPath)
			}

			hasher, err := securityContext.NewHasher()
			if err != nil {
				return err
			}

			if _, err = hasher.Write(blockData[:block.Size]); err != nil {
				return err
			}
			hash := fmt.Sprintf("%x", hasher.Sum(nil))

			if hash != block.Hash {
				log.WithFields(
					log.Fields{
						"actual_hash":   hash,
						"expected_hash": block.Hash}).Error("corrupt block detected")

				return fmt.Errorf("corrupt block: %s", blockPath)
			}

			if !dryRun {
				if _, err = outputFile.Write(blockData[:block.Size]); err != nil {
					log.WithError(err).Error("failed to restore block")
					return err
				}
			}

			return nil
		}

		if err := restoreBlock(); err != nil {
			return err
		}
	}

	log.WithFields(
		log.Fields{
			"path":          file.Path,
			"restored_path": outputPath}).Debug("file restored")

	return nil
}
