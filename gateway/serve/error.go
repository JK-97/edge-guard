package serve

import (
	"encoding/json"
	// "log"
	"net/http"

	"jxcore/gateway/log"
)

// APIError api 异常类型
type APIError struct {
	Code        int    `json:"code"`
	Description string `json:"desc"`
}

const mimeJSON = "application/json; charset=utf-8"

// Error 服务报错
func Error(w http.ResponseWriter, reason string, code int) {
	w.Header().Set("Content-Type", mimeJSON)

	e := APIError{code, reason}

	w.WriteHeader(code)

	rs, err := json.Marshal(e)
	if err != nil {
		log.Fatalln(err)
	}

	w.Write(rs)
}

// ErrorWithCode 使用预设的 Http 状态码抛出异常
func ErrorWithCode(w http.ResponseWriter, code int) {
	Error(w, http.StatusText(code), code)
}

// ErrorNotFound  服务报错
func ErrorNotFound(w http.ResponseWriter) {
	ErrorWithCode(w, http.StatusNotFound)
}

type notFoundHandler struct {
}

func (h *notFoundHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	log.Printf("In: 404 \t%s %s %s\n", r.RemoteAddr, r.Method, r.URL)

	ErrorNotFound(w)
}

// NewNotFoundHandler 404 处理问题
func NewNotFoundHandler() http.Handler {
	return &notFoundHandler{}
}
