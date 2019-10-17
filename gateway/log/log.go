package log

import (
	"log"
	"os"
	"time"

	"github.com/sirupsen/logrus"
)

// These flags define which text to prefix to each log entry generated by the Logger.
// Bits are or'ed together to control what's printed.
// There is no control over the order they appear (the order listed
// here) or the format they present (as described in the comments).
// The prefix is followed by a colon only when Llongfile or Lshortfile
// is specified.
// For example, flags Ldate | Ltime (or LstdFlags) produce,
//	2009/01/23 01:23:23 message
// while flags Ldate | Ltime | Lmicroseconds | Llongfile produce,
//	2009/01/23 01:23:23.123123 /a/b/c/d.go:23: message
const (
	Ldate         int = log.Ldate
	Ltime         int = log.Ltime
	Lmicroseconds int = log.Lmicroseconds
	Llongfile     int = log.Llongfile
	Lshortfile    int = log.Lshortfile
	LUTC          int = log.LUTC
	LstdFlags     int = log.LstdFlags
)

type logHook interface {
	logrus.Hook
	SetReportCaller(bool)
}

var (
	devLogger  logrus.Ext1FieldLogger
	prodLogger logrus.Ext1FieldLogger
	formatter  logrus.Formatter

	hook logHook
)

type callerHook struct {
	levels       []logrus.Level
	ReportCaller bool
}

func (h *callerHook) Levels() []logrus.Level {
	return h.levels
}

func (h *callerHook) Fire(entry *logrus.Entry) error {
	if h.ReportCaller {
		caller := getCaller()
		if caller != nil {
			entry.Data[logrus.FieldKeyFile] = caller.File
			entry.Data[logrus.FieldKeyFunc] = caller.Function
		}
	}

	return nil
}

func (h *callerHook) SetReportCaller(include bool) {
	h.ReportCaller = include
}

func init() {
	formatter = &logrus.JSONFormatter{
		TimestampFormat: time.RFC3339,
	}

	devLogger = logrus.New()
	prodLogger = logrus.StandardLogger()

	// logrus.SetReportCaller(true)
	devLogger.(*logrus.Logger).Level = logrus.TraceLevel
	devLogger.(*logrus.Logger).Out = os.Stdout
	devLogger.(*logrus.Logger).SetFormatter(formatter)
	prodLogger.(*logrus.Logger).SetFormatter(formatter)

	hook = &callerHook{
		levels: logrus.AllLevels,
	}

	devLogger.(*logrus.Logger).AddHook(hook)
	// prodLogger.(*logrus.Logger).AddHook(hook)
}

// SetReportCaller sets whether the standard logger will include the calling
// method as a field.
func SetReportCaller(include bool) {
	hook.SetReportCaller(include)
}

// Flags returns the output flags for the standard logger.
func Flags() int {
	return log.Flags()
}

// SetFlags sets the output flags for the standard logger.
func SetFlags(flag int) {
	log.SetFlags(flag)
}

// Trace logs a message at level Trace on the standard logger.
func Trace(args ...interface{}) {
	devLogger.Trace(args...)
}

// Traceln logs a message at level Trace on the standard logger.
func Traceln(args ...interface{}) {
	devLogger.Traceln(args...)
}

// Tracef logs a message at level Trace on the standard logger.
func Tracef(format string, args ...interface{}) {
	devLogger.Tracef(format, args...)
}

// Debug logs a message at level Debug on the standard logger.
func Debug(args ...interface{}) {
	devLogger.Debug(args...)
}

// Debugln logs a message at level Debug on the standard logger.
func Debugln(args ...interface{}) {
	devLogger.Debugln(args...)
}

// Debugf logs a message at level Debug on the standard logger.
func Debugf(format string, args ...interface{}) {
	devLogger.Debugf(format, args...)
}

// Print logs a message at level Info on the standard logger.
func Print(args ...interface{}) {
	prodLogger.Print(args...)
}

// Println logs a message at level Info on the standard logger.
func Println(args ...interface{}) {
	prodLogger.Println(args...)
}

// Printf logs a message at level Info on the standard logger.
func Printf(format string, args ...interface{}) {
	prodLogger.Printf(format, args...)
}

// Info logs a message at level Info on the standard logger.
func Info(args ...interface{}) {
	prodLogger.Info(args...)
}

// Infoln logs a message at level Info on the standard logger.
func Infoln(args ...interface{}) {
	prodLogger.Infoln(args...)
}

// Infof logs a message at level Info on the standard logger.
func Infof(format string, args ...interface{}) {
	prodLogger.Infof(format, args...)
}

// Warning logs a message at level Warning on the standard logger.
func Warning(args ...interface{}) {
	prodLogger.Warning(args...)
}

// Warningln logs a message at level Warning on the standard logger.
func Warningln(args ...interface{}) {
	prodLogger.Warningln(args...)
}

// Warningf logs a message at level Warning on the standard logger.
func Warningf(format string, args ...interface{}) {
	prodLogger.Warningf(format, args...)
}

// Error logs a message at level Error on the standard logger.
func Error(args ...interface{}) {
	prodLogger.Error(args...)
}

// Errorln logs a message at level Error on the standard logger.
func Errorln(args ...interface{}) {
	prodLogger.Errorln(args...)
}

// Errorf logs a message at level Error on the standard logger.
func Errorf(format string, args ...interface{}) {
	prodLogger.Errorf(format, args...)
}

// Fatal logs a message at level Fata on the standard logger.
func Fatal(args ...interface{}) {
	devLogger.Fatal(args...)
}

// Fatalln logs a message at level Fata on the standard logger.
func Fatalln(args ...interface{}) {
	devLogger.Fatalln(args...)
}

// Fatalf logs a message at level Fata on the standard logger.
func Fatalf(format string, args ...interface{}) {
	devLogger.Fatalf(format, args...)
}

// Panic logs a message at level Panic on the standard logger.
func Panic(args ...interface{}) {
	devLogger.Panic(args...)
}

// Panicf logs a message at level Panic on the standard logger.
func Panicf(format string, args ...interface{}) {
	devLogger.Panicf(format, args...)
}

// Panicln logs a message at level Panic on the standard logger.
func Panicln(args ...interface{}) {
	devLogger.Panicln(args...)
}