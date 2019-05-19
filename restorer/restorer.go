package restorer

import (
	"github.com/mboye/kopi/input"
	"github.com/mboye/kopi/security"
	"github.com/mboye/kopi/stage"

	_ "github.com/mboye/kopi/loglevel"
	"github.com/mboye/kopi/model"
	log "github.com/sirupsen/logrus"
)

type restorer struct {
	inputDir, outputDir string
	dryRun              bool
	decrypt             bool
	progressInterval    int
}

var _ stage.Stage = (*restorer)(nil)

func New(inputDir, outputDir string, dryRun, decrypt bool, progressInterval int) (stage.Stage, error) {
	return &restorer{inputDir, outputDir, dryRun, decrypt, progressInterval}, nil
}

func (r *restorer) Execute() error {
	securityContext, err := security.NewContext(r.inputDir, r.decrypt)
	if err != nil {
		log.WithField("error", err).Fatal("failed to create security context")
	}

	log.WithField("destination", r.outputDir).Info("beginning to restore files")

	restoreFile := func(file *model.File) (err error) {
		if file.Mode.IsDir() {
			return restoreDir(file, r.outputDir, r.dryRun)
		} else {
			return restoreFile(file, r.inputDir, r.outputDir, securityContext, r.dryRun)
		}
	}

	return input.ProcessFilesWithProgress(restoreFile, uint(r.progressInterval))
}
