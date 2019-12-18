package controller

import (
	"io/ioutil"
	"net/http"

	log "jxcore/lowapi/logger"

	"encoding/json"
)

type bootstrapRequest struct {
	RegistrationCode string `json:"registercode"`
	DhcpServer       string `json:"dhcpserver"`
}

func Boostrap(w http.ResponseWriter, r *http.Request) {
	data, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Error(err)
	}
	request := bootstrapRequest{}
	err = json.Unmarshal(data, &request)
	if err != nil {
		panic(err)
	}
	if request.RegistrationCode == "" && request.DhcpServer == "" {
		return
	}

}
