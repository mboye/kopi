package util

import (
	"os"

	log "github.com/sirupsen/logrus"
)

// SetLogLevel detects and sets the log level from the environment variable KOPI_LOG_LEVEL.
// The following values are supported: DEBUG, INFO, WARN, ERROR, FATAL, and PANIC.
func SetLogLevel() {
	if level, exists := os.LookupEnv("KOPI_LOG_LEVEL"); exists {
		switch level {
		case "DEBUG":
			log.SetLevel(log.DebugLevel)
		case "INFO":
			log.SetLevel(log.InfoLevel)
		case "WARN":
			log.SetLevel(log.WarnLevel)
		case "ERROR":
			log.SetLevel(log.ErrorLevel)
		case "FATAL":
			log.SetLevel(log.FatalLevel)
		case "PANIC":
			log.SetLevel(log.PanicLevel)
		default:
			log.Warnf("Unknown log level: %s", level)
		}
	}
}
