package logger

import (
	"jxcore/lowapi/data"
	"jxcore/version"

	"gitlab.jiangxingai.com/applications/base-modules/internal-sdk/go-utils/logger"
)

func init() {
	logger.SetLoggerConfig(logger.Configuration{
		CallerSkip: 1,
	})
}

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
func Printf(format string, args ...interface{}) {
	logger.WithFields(logger.Fields{"JXCORE_VERSION": version.Version}).Infof(format, args...)
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

func Println(args ...interface{}) {
	logger.WithFields(logger.Fields{"JXCORE_VERSION": version.Version}).Info(args...)
}
func Fatalln(args ...interface{}) {
	logger.WithFields(logger.Fields{"JXCORE_VERSION": version.Version}).Fatal(args...)
}

func WithFields(keyValues Fields) logger.Logger {
	cp := data.CopyMap(keyValues)
	cp["JXCORE_VERSION"] = version.Version
	return logger.WithFields(logger.Fields(cp))
}
