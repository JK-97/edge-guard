package version

import (
	"fmt"
	"runtime"
)

// GitCommit The git commit that was compiled. This will be filled in by the compiler.
var GitCommit string

// Version The main version number that is being run at the moment.
var Version = "1.0.0"

// BuildDate 编译日期
var BuildDate = "2019.08.23"

// GoVersion Go version
var GoVersion = runtime.Version()

// OsArch 架构
var OsArch = fmt.Sprintf("%s %s", runtime.GOOS, runtime.GOARCH)
