package controller

import "net/http"

import "io/ioutil"

import "jxcore/gateway/log"

import "encoding/json"

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
