package controller

import (
	"encoding/json"
	"html/template"
	"io/ioutil"
	log "jxcore/lowapi/logger"
	"net/http"

	"gitlab.jiangxingai.com/applications/base-modules/internal-sdk/go-utils/logger"
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
	w.WriteHeader(statusCode)
	w.Header().Set("Content-Type", "application/json")
	_, err = w.Write(b)
	if err != nil {
		logger.Error(err)
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

	w.Write(body)
	return nil
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
