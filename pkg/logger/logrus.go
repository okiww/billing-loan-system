package logger

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"os"
)

var log *logrus.Logger

// InitLogger initializes the logrus logger with a specified log level and log output
func InitLogger(logLevel string, logFile string) {
	log = logrus.New()

	// Set log level
	level, err := logrus.ParseLevel(logLevel)
	if err != nil {
		log.Fatalf("Invalid log level: %v", err)
	}
	log.SetLevel(level)

	// Set log output (console or file)
	if logFile != "" {
		file, err := os.OpenFile(logFile, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
		if err != nil {
			log.Fatalf("Error opening log file: %v", err)
		}
		log.SetOutput(file)
	} else {
		// Default output is stderr (console)
		log.SetOutput(os.Stdout)
	}
}

// GetLogger returns the logger instance
func GetLogger() *logrus.Logger {
	if log == nil {
		log = logrus.New()
	}
	return log
}

func Info(msg string) {
	GetLogger().Info()
}

func Infof(format string, v ...any) {
	GetLogger().Info(fmt.Sprintf(format, v...))
}

func Fatal(msg string) {
	GetLogger().Fatal(msg)
}

func Fatalf(format string, v ...any) {
	Fatal(fmt.Sprintf(format, v...))
}
