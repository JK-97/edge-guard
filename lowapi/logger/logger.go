package logger

import (
	"jxcore/version"

	"gitlab.jiangxingai.com/applications/base-modules/internal-sdk/go-utils/logger"
)

type Fields logger.Fields

func Debugf(format string, args ...interface{}) {
	logger.WithFields(logger.Fields{"JXCORE_VERSION": version.Version}).Debugf(format, args...)
}

func Infof(format string, args ...interface{}) {
	logger.WithFields(logger.Fields{"JXCORE_VERSION": version.Version}).Infof(format, args...)
}

func Warnf(format string, args ...interface{}) {
	logger.WithFields(logger.Fields{"JXCORE_VERSION": version.Version}).Warnf(format, args...)
}

func Errorf(format string, args ...interface{}) {
	logger.WithFields(logger.Fields{"JXCORE_VERSION": version.Version}).Errorf(format, args...)
}

func Fatalf(format string, args ...interface{}) {
	logger.WithFields(logger.Fields{"JXCORE_VERSION": version.Version}).Fatalf(format, args...)
}

func Panicf(format string, args ...interface{}) {
	logger.WithFields(logger.Fields{"JXCORE_VERSION": version.Version}).Panicf(format, args...)
}

func Debug(args ...interface{}) {
	logger.WithFields(logger.Fields{"JXCORE_VERSION": version.Version}).Debug(args...)
}

func Info(args ...interface{}) {
	logger.WithFields(logger.Fields{"JXCORE_VERSION": version.Version}).Info(args...)
}

func Warn(args ...interface{}) {
	logger.WithFields(logger.Fields{"JXCORE_VERSION": version.Version}).Warn(args...)
}

func Error(args ...interface{}) {
	logger.WithFields(logger.Fields{"JXCORE_VERSION": version.Version}).Error(args...)
}

func Fatal(args ...interface{}) {
	logger.WithFields(logger.Fields{"JXCORE_VERSION": version.Version}).Fatal(args...)
}

func Panic(args ...interface{}) {
	logger.WithFields(logger.Fields{"JXCORE_VERSION": version.Version}).Panic(args...)
}

func WithFields(keyValues Fields) logger.Logger {
	keyValues["JXCORE_VERSION"] = version.Version
	return logger.WithFields(logger.Fields(keyValues))
}
