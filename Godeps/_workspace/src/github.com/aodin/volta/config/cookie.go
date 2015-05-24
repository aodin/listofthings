package config

import (
	"time"
)

// CookieConfig contains the fields needed to set and retrieve cookies.
// Cookie names are valid tokens as defined by RFC 2616 section 2.2:
// http://tools.ietf.org/html/rfc2616#section-2.2
// TL;DR: Any non-control or non-separator character.
type CookieConfig struct {
	Age      time.Duration `json:"age"`
	Domain   string        `json:"domain"`
	HttpOnly bool          `json:"http_only"`
	Name     string        `json:"name"`
	Path     string        `json:"path"`
	Secure   bool          `json:"secure"`
}

// DefaultCookie is a default CookieConfig implementation. It expires after
// two weeks and is not very secure.
var DefaultCookie = CookieConfig{
	Age:      14 * 24 * time.Hour, // Two weeks
	Domain:   "",
	HttpOnly: false,
	Name:     "sessionid",
	Path:     "/",
	Secure:   false,
}
