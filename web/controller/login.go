package controller

import (
	"fmt"
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
	unmarshalJson(r.Body, &request)

	password, err := getPassword()
	if err != nil {
		panic(err)
	}
	if password != request.Password {
		panic(HTTPError{
			Code: http.StatusUnauthorized,
			Err:  fmt.Errorf("password incorrect."),
		})
	}

	user := User{
		Authenticated: true,
	}
	if err := saveSessionUser(w, r, user, false); err != nil {
		panic(err)
	}
}

func PostLogout(w http.ResponseWriter, r *http.Request) {
	if err := saveSessionUser(w, r, User{}, true); err != nil {
		panic(err)
	}
}

func GetSessionUser(r *http.Request) (User, error) {
	session, _ := getSession(r)
	val := session.Values[sessionKeyUser]
	user, ok := val.(User)
	if !ok {
		return User{Authenticated: false}, nil
	}
	return user, nil
}

func saveSessionUser(w http.ResponseWriter, r *http.Request, user User, expires bool) error {
	session, _ := getSession(r)
	session.Values[sessionKeyUser] = user
	if expires {
		session.Options.MaxAge = -1
	}
	return session.Save(r, w)
}
