package version

import (
	"fmt"
	"runtime"
	"strings"
)

var (
	GitCommit = "unknown"
	Version   = "unknown"
	BuildDate = "unknown"
	GitStatus = "unknown"
	GoVersion = runtime.Version()
	OsArch    = fmt.Sprintf("%s %s", runtime.GOOS, runtime.GOARCH)
)

// 返回单行格式
func StringifySingleLine() string {
	return fmt.Sprintf("GitCommitLog=%s. GitStatus=%s. BuildTime=%s. GoVersion=%s. runtime=%s.",
		GitCommit, GitStatus, BuildDate, GoVersion, OsArch)
}

// 返回多行格式
func StringifyMultiLine() string {
	return fmt.Sprintf("GitCommitLog=%s\nGitStatus=%s\nBuildTime=%s\nGoVersion=%s\nruntime=%s\n",
		GitCommit, GitStatus, BuildDate, GoVersion, OsArch)
}

// 对一些值做美化处理
func beauty() {
	if GitStatus == "" {
		// GitStatus 为空时，说明本地源码与最近的 commit 记录一致，无修改
		// 为它赋一个特殊值
		GitStatus = "cleanly"
	} else {
		// 将多行结果合并为一行
		GitStatus = strings.Replace(strings.Replace(GitStatus, "\r\n", " |", -1), "\n", " |", -1)
	}
}

func init() {
	beauty()
}
