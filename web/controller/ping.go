package controller

import (
	"net/http"
)

func Ping(w http.ResponseWriter, r *http.Request) {
	RespondJSON("pong", w, 200)
}
