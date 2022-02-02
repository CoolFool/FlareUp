package logging

import (
	"github.com/sirupsen/logrus"
	"os"
)

type Log struct {
	Logger *logrus.Logger
}

func Init() Log {
	formatter := &logrus.TextFormatter{
		FullTimestamp: true,
		DisableColors: false,
	}
	log := Log{Logger: logrus.New()}
	log.Logger.SetFormatter(formatter)
	log.Logger.SetOutput(os.Stdout)
	log.Logger.SetLevel(logrus.InfoLevel)
	return log
}

func (l *Log) Error(any ...interface{}) {
	l.Logger.Error(any)
}

func (l *Log) Info(any ...interface{}) {
	l.Logger.Info(any)
}

func (l *Log) Warn(any ...interface{}) {
	l.Logger.Warn(any)
}

func (l *Log) Debug(any ...interface{}) {
	l.Logger.Debug(any)
}

func (l *Log) Fatal(any ...interface{}) {
	l.Logger.Fatal(any)
}
