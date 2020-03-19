package system

import (
	"encoding/json"
	"io/ioutil"
	"net/http"

	"gitlab.jiangxingai.com/applications/base-modules/internal-sdk/go-utils/logger"
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
