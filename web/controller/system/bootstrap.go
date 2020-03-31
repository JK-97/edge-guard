package system

import (
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/JK-97/go-utils/logger"
)

type bootstrapRequest struct {
	RegistrationCode string `json:"registercode"`
	DhcpServer       string `json:"dhcpserver"`
}

func Boostrap(w http.ResponseWriter, r *http.Request) {
	data, err := ioutil.ReadAll(r.Body)
	if err != nil {
		logger.Error(err)
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
