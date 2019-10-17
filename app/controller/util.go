package controller

import (
	"encoding/json"
	"html/template"
	"io/ioutil"
	"jxcore/app/schema"
	"jxcore/log"
	"net/http"
)


// Error handles server error
func Error(w http.ResponseWriter, err error, code int) {
	if h, ok := err.(schema.HTTPError); ok {
		code = h.Code
	}

	http.Error(w, http.StatusText(code), code)
	log.Error(err)
}

func respondSuccessJSON(obj interface{}, w http.ResponseWriter, r *http.Request) {
	Resp:= schema.BaseResp{Data:obj,Desc:"success"}
	b, err := json.Marshal(Resp)
	if err != nil {
		Error(w, err, http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(b)
}

func respondResonJSON(obj interface{}, w http.ResponseWriter, r *http.Request,reson string) {

	Resp:= schema.BaseResp{Data:obj,Desc:reson}
	b, err := json.Marshal(Resp)
	if err != nil {
		Error(w, err, http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(b)
}



func serveStatic(path string, w http.ResponseWriter, r *http.Request) error {
	body, err := ioutil.ReadFile(path)
	if err != nil {
		Error(w, err, http.StatusInternalServerError)
		return err
	}

	w.Write(body)
	return nil
}

func serveTemplate(path string, data interface{}, w http.ResponseWriter, r *http.Request) error {
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