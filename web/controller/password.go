package controller

import (
	"fmt"
	"jxcore/lowapi/store/filestore"
	"net/http"
)

type SetPasswordRequest struct {
	OldPassword string `json:"old_password"`
	NewPassword string `json:"new_password"`
}

const (
	passwordKey     = "password"
	defaultPassword = "d04db14aa76e3a03c4e383136f941c0d" // jiangxing123
)

// 设置Jxcore密码
func SetPasswordHandler(w http.ResponseWriter, r *http.Request) {
	request := SetPasswordRequest{}
	unmarshalJson(r.Body, &request)

	oldPassword, err := getPassword()
	if err != nil {
		panic(err)
	}
	if oldPassword != request.OldPassword {
		panic(HTTPError{
			Code: http.StatusBadRequest,
			Err:  fmt.Errorf("old password incorrect."),
		})
	}

	err = filestore.KV.Set(passwordKey, []byte(request.NewPassword))
	if err != nil {
		panic(err)
	}
	RespondSuccessJSON(nil, w)
}

func getPassword() (string, error) {
	data, err := filestore.KV.GetDefault(passwordKey, []byte(defaultPassword))
	return string(data), err
}
