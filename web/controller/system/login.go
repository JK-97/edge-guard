package system

import (
	"fmt"
	"jxcore/oplog"
	"jxcore/oplog/logs"
	"jxcore/oplog/types"
	"jxcore/web/controller/utils"
	"net/http"
)

const (
	sessionKeyUser = "user"
)

type loginReq struct {
	Password string `json:"password"`
}

func PostLogin(w http.ResponseWriter, r *http.Request) {
	request := loginReq{}
	utils.MustUnmarshalJson(r.Body, &request)

	password, err := getPassword()
	if err != nil {
		panic(err)
	}
	if password != request.Password {
		panic(utils.HTTPError{
			Code: http.StatusUnauthorized,
			Err:  fmt.Errorf("password incorrect."),
		})
	}

	user := utils.User{
		Authenticated: true,
	}
	if err := saveSessionUser(w, r, user, false); err != nil {
		panic(err)
	}
	oplog.Insert(logs.NewOplog(types.AUTH, "somce one login"))
}

func PostLogout(w http.ResponseWriter, r *http.Request) {
	if err := saveSessionUser(w, r, utils.User{}, true); err != nil {
		panic(err)
	}
	oplog.Insert(logs.NewOplog(types.AUTH, "somce one logout"))
}

func GetSessionUser(r *http.Request) (utils.User, error) {
	session, _ := utils.GetLoginSession(r)
	val := session.Values[sessionKeyUser]
	user, ok := val.(utils.User)
	if !ok {
		return utils.User{Authenticated: false}, nil
	}
	return user, nil
}

func saveSessionUser(w http.ResponseWriter, r *http.Request, user utils.User, expires bool) error {
	session, _ := utils.GetLoginSession(r)
	session.Values[sessionKeyUser] = user
	if expires {
		session.Options.MaxAge = -1
	}
	return session.Save(r, w)
}
