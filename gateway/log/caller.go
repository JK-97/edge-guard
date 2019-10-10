package log

import (
	"path/filepath"
	"runtime"
	"strings"
	"sync"
)

var (
	// bufferPool *sync.Pool

	// qualified package name, cached at first use
	logPackage string

	rootPackage string

	// qualified package name, cached at first use
	logrusPackage = "github.com/sirupsen/logrus"

	// Positions in the call stack when tracing to report the calling method
	minimumCallerDepth int

	// Used for caller information initialisation
	callerInitOnce sync.Once
)

const (
	maximumCallerDepth int = 25
	knownLogrusFrames  int = 6
)

// getPackageName reduces a fully qualified function name to the package name
// There really ought to be to be a better way...
func getPackageName(f string) string {
	for {
		lastPeriod := strings.LastIndex(f, ".")
		lastSlash := strings.LastIndex(f, "/")
		if lastPeriod > lastSlash {
			f = f[:lastPeriod]
		} else {
			break
		}
	}

	return f
}

// getCaller retrieves the name of the first non-logrus calling function
func getCaller() *runtime.Frame {

	// cache this package's fully-qualified name
	callerInitOnce.Do(func() {
		pcs := make([]uintptr, 2)
		_ = runtime.Callers(0, pcs)
		logPackage = getPackageName(runtime.FuncForPC(pcs[1]).Name())
		rootPackage = strings.ReplaceAll(filepath.Dir(logPackage), "\\", "/")

		// now that we have the cache, we can skip a minimum count of known-logrus functions
		// XXX this is dubious, the number of frames may vary
		minimumCallerDepth = knownLogrusFrames
	})

	// Restrict the lookback frames to avoid runaway lookups
	pcs := make([]uintptr, maximumCallerDepth)
	depth := runtime.Callers(minimumCallerDepth, pcs)
	frames := runtime.CallersFrames(pcs[:depth])

	for f, again := frames.Next(); again; f, again = frames.Next() {
		pkg := getPackageName(f.Function)

		// If the caller isn't part of this package, we're done
		if pkg != logrusPackage && pkg != logPackage {
			trimFilePath(&f.File)
			return &f
		}
	}

	// if we got here, we failed to find the caller's context
	return nil
}

func trimFilePath(f *string) {

	i := strings.Index(*f, "src")
	if i > 0 {
		if j := strings.Index(*f, rootPackage); j > 0 {
			*f = (*f)[j+len(rootPackage)+1:]
		} else {
			*f = (*f)[i+4:]
		}
	}
}
