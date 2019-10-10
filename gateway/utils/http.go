package utils

import (
	"encoding/json"
	"net/http"
)

func WriteJson(w http.ResponseWriter, obj interface{}) {
	js, err := json.Marshal(obj)
	if err != nil {
		panic(err)
	}
	w.Header().Set("Content-Type", "application/json")
	_, err = w.Write(js)
	if err != nil {
		panic(err)
	}
}
