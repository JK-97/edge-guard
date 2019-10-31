#!/usr/bin/env bash
# 参数:  版本号   编译平台
# 编译平台缺省时,为当前平台
set -x
Version=${1}
BuildArch="${2}"
if [ "${BuildArch}" == "" ]; then
    BuildArch=$(arch)
fi
# 获取源码最近一次 git commit log，包含 commit sha 值，以及 commit message
GitCommitLog=$(git log --pretty=oneline -n 1)
# 将 log 原始字符串中的单引号替换成双引号
GitCommitLog=${GitCommitLog//\'/\"}
# 检查源码在git commit 基础上，是否有本地修改，且未提交的内容
GitStatus=$(git status -s)
# 获取当前时间
BuildTime=$(date +'%Y.%m.%d-%H:%M:%S')
# 获取 Go 的版本
# BuildGoVersion=`go version`

LDFlags=" \
    -X 'jxcore/version.GitCommit=${GitCommitLog}' \
    -X 'jxcore/version.GitStatus=${GitStatus}' \
    -X 'jxcore/version.BuildDate=${BuildTime}' \
    -X 'jxcore/version.Version=${Version}' \
"
echo ${Version}

if [ $BuildArch == "arm64" ]; then
    CGO_ENABLED=0 GOOS=linux GOARCH=arm64 go build -ldflags "$LDFlags"
else
    go build -ldflags "$LDFlags"
fi

echo "build ${BuildArch} done."