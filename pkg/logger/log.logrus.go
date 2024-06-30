package logger

import (
	"io"
	"os"

	"github.com/james-wukong/go-app/internal/constants"
	"github.com/sirupsen/logrus"
)

var log *logrus.Logger

func init() {
	// create a new logger instance
	log = logrus.New()
	logFile := constants.LoggerPathFile
	// log_file := config.AppConfig.LogFile
	fileOut, err := os.OpenFile(logFile, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		log.Fatal("Failed to open log file:", logFile, err)
	}
	// Create a MultiWriter with os.Stdout and the log file
	mw := io.MultiWriter(os.Stdout, fileOut)
	// mw := io.MultiWriter(os.Stdout)

	// set output to channels
	log.SetOutput(mw)
	// Only log the debug severity or above.
	log.SetLevel(logrus.DebugLevel)
	// add the calling method as a field
	log.SetReportCaller(true)

	// set log formatter
	log.SetFormatter(&logrus.TextFormatter{
		ForceColors:     true,
		FullTimestamp:   true,
		TimestampFormat: "2006-01-02 15:04:05",
	})

	// default fields
	log.WithFields(logrus.Fields{constants.LoggerCategory: constants.LoggerCategoryServer})
}

func Info(message string, fields logrus.Fields) {
	log.WithFields(fields).Info(message)
}

func InfoF(format string, fields logrus.Fields, args ...interface{}) {
	log.WithFields(fields).Infof(format, args...)
}

func Debug(message string, fields logrus.Fields) {
	log.WithFields(fields).Debug(message)
}

func DebugF(format string, fields logrus.Fields, args ...interface{}) {
	log.WithFields(fields).Debugf(format, args...)
}

func Error(message string, fields logrus.Fields) {
	log.WithFields(fields).Error(message)
}

func ErrorF(format string, fields logrus.Fields, args ...interface{}) {
	log.WithFields(fields).Errorf(format, args...)
}

func Fatal(message string, fields logrus.Fields) {
	log.WithFields(fields).Fatal(message)
}

func FatalF(format string, fields logrus.Fields, args ...interface{}) {
	log.WithFields(fields).Fatalf(format, args...)
}

func Panic(message string, fields logrus.Fields) {
	log.WithFields(fields).Panic(message)
}

func PanicF(format string, fields logrus.Fields, args ...interface{}) {
	log.WithFields(fields).Panicf(format, args...)
}
