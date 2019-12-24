package heartbeat

import (
	"errors"
	"io"
	"net"
	"strings"
)

// isConnectionBroken 判断连接是否仍然可用
func isConnectionBroken(err error) bool {
	if err == io.EOF {
		return true
	}
	for _, aErr := range AbortErrs {
		if errors.Is(err, aErr) {
			return true
		}
	}

	return strings.Contains(err.Error(), "closed network")
}

// isDNSError 判断是否是 DNS 解析错误
func isDNSError(err error) bool {
	var e error

	for e = errors.Unwrap(err); e != nil; e = errors.Unwrap(e) {
		if _, ok := e.(*net.DNSError); ok {
			return true
		}
	}
	return false
}

// isDialErr 判断是否是因为域名解析错误导致的错误
func isDialErr(err error) bool {
	for _, aErr := range AbortErrs {
		if errors.Is(err, aErr) {
			return true
		}
	}
	return false
}
