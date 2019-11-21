package controller

import (
	"net/http"
)

func Ping(w http.ResponseWriter, r *http.Request) {
	respondJSON("pong", w, r)

}
