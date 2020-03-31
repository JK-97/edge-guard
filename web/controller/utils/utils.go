package utils

import (
	"encoding/json"
	log "github.com/JK-97/edge-guard/lowapi/logger"
	"html/template"
	"io"
	"io/ioutil"
	"net/http"
	"strconv"

	"github.com/JK-97/go-utils/logger"
)

// Error handles server error
func Error(w http.ResponseWriter, err error, code int) {
	if h, ok := err.(HTTPError); ok {
		code = h.Code
	}

	RespondReasonJSON(nil, w, err.Error(), code)
	log.Error(err)
}

func RespondJSON(obj interface{}, w http.ResponseWriter, statusCode int) {
	if obj == nil {
		return
	}

	b, err := json.Marshal(obj)
	if err != nil {
		Error(w, err, http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	_, err = w.Write(b)
	if err != nil {
		logger.Error(err)
	}
	if f, ok := w.(http.Flusher); ok {
		f.Flush()
	}
}

func RespondReasonJSON(obj interface{}, w http.ResponseWriter, reason string, statusCode int) {
	Resp := BaseResp{Data: obj, Desc: reason}
	RespondJSON(Resp, w, statusCode)
}

func RespondSuccessJSON(obj interface{}, w http.ResponseWriter) {
	RespondReasonJSON(obj, w, "success", 200)
}

func ServeStatic(path string, w http.ResponseWriter, r *http.Request) error {
	body, err := ioutil.ReadFile(path)
	if err != nil {
		Error(w, err, http.StatusInternalServerError)
		return err
	}

	_, err = w.Write(body)
	return err
}

func ServeTemplate(path string, data interface{}, w http.ResponseWriter, r *http.Request) error {
	body, err := ioutil.ReadFile(path)
	if err != nil {
		Error(w, err, http.StatusInternalServerError)
		return err
	}

	t, err := template.New(path).Parse(string(body))
	if err != nil {
		Error(w, err, http.StatusInternalServerError)
		return err
	}

	err = t.Execute(w, data)
	if err != nil {
		Error(w, err, http.StatusInternalServerError)
		return err
	}

	return nil
}

func CatchPanic(w http.ResponseWriter, r *http.Request, statusCode int) {
	if err := recover(); err != nil {
		log.Error("handler failed", err)
		RespondJSON(err, w, statusCode)
	}
}

func MustUnmarshalJson(body io.ReadCloser, out interface{}) {
	err := UnmarshalJson(body, out)
	if err != nil {
		panic(err)
	}
}

func UnmarshalJson(body io.ReadCloser, out interface{}) error {
	data, err := ioutil.ReadAll(body)
	if err != nil {
		return err
	}
	return json.Unmarshal(data, out)
}

func GetPageInfo(r *http.Request) (int, int) {
	offset := 0
	v := r.URL.Query()["offset"]
	if len(v) == 1 {
		offset, _ = strconv.Atoi(v[0])
	}
	limit := 10
	v = r.URL.Query()["page"]
	if len(v) == 0 {
		limit = 1 << 16
	}

	return offset, limit
}
