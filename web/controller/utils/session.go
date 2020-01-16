package utils

import (
	"encoding/gob"
	"net/http"

	"github.com/gorilla/securecookie"
	"github.com/gorilla/sessions"
)

type User struct {
	Authenticated bool
}

var store *sessions.CookieStore

func init() {
	authKeyOne := securecookie.GenerateRandomKey(64)
	encryptionKeyOne := securecookie.GenerateRandomKey(32)

	store = sessions.NewCookieStore(
		authKeyOne,
		encryptionKeyOne,
	)

	store.Options = &sessions.Options{
		MaxAge:   60 * 60 * 24,
		HttpOnly: true,
	}

	gob.Register(User{})
}

func GetLoginSession(r *http.Request) (*sessions.Session, error) {
	return store.Get(r, "jiangxing_login_jxcore_ticket")
}
