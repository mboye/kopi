package restorer

import (
	"fmt"
	"os"

	"github.com/mboye/kopi/model"
	log "github.com/sirupsen/logrus"
)

func restoreDir(file *model.File, outputDir string, dryRun bool) error {
	outputPath := fmt.Sprintf("%s/%s", outputDir, file.Path)
	log.WithFields(log.Fields{
		"path": outputPath,
		"mode": file.Mode}).Debug("restoring directory")

	if dryRun {
		return nil
	}
	return os.MkdirAll(outputPath, file.Mode)
}
