package controller

import (
	"jxcore/internal/config"
	"net/http"
)

func Reset(w http.ResponseWriter, r *http.Request) {
	err := config.ResetSystemConfig()
	if err != nil {
		panic(err)
	}
	RespondSuccessJSON(nil, w)
}
