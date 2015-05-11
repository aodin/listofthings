package auth

import (
	"net/http"

	"github.com/aodin/volta/config"
)

func SetCookie(w http.ResponseWriter, c config.CookieConfig, session Session) {
	cookie := &http.Cookie{
		Name:     c.Name,
		Value:    session.Key,
		Path:     c.Path,
		Domain:   c.Domain,
		Expires:  session.Expires,
		HttpOnly: c.HttpOnly,
		Secure:   c.Secure,
	}
	http.SetCookie(w, cookie)
}
