package logging

import (
	"github.com/sirupsen/logrus"
)

var logger *logrus.Logger

func InitLog(debug bool) {
	var level logrus.Level
	if debug {
		level = logrus.DebugLevel
	} else {
		level = logrus.WarnLevel
	}

	logger = logrus.New()
	logger.SetLevel(level)
	logger.SetFormatter(&logrus.JSONFormatter{
		TimestampFormat: "2006/01/02 15:04:05",
	})
}

func Logger() *logrus.Entry {
	return logrus.NewEntry(logger)
}
