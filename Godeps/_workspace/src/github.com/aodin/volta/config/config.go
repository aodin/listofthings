package config

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/url"
	"os"
)

// Config is the parent configuration struct and includes fields for single
// configurations of a database, cookie, and SMTP connection.
type Config struct {
	HTTPS       bool           `json:"https"`
	Domain      string         `json:"domain"`
	ProxyDomain string         `json:"proxy_domain"`
	Port        int            `json:"port"`
	ProxyPort   int            `json:"proxy_port"`
	TemplateDir string         `json:"templates"`
	AbsPath     string         `json:"abs_path"`
	MediaDir    string         `json:"media"`
	MediaURL    string         `json:"media_url"`
	StaticDir   string         `json:"static"`
	StaticURL   string         `json:"static_url"`
	SecretKey   string         `json:"secret_key"`
	Database    DatabaseConfig `json:"database"`
	Cookie      CookieConfig   `json:"cookie"`
	SMTP        SMTPConfig     `json:"smtp"`
	Metadata    Metadata       `json:"metadata"`
}

// Address returns the domain:port pair.
func (c Config) Address() string {
	return fmt.Sprintf("%s:%d", c.Domain, c.Port)
}

// URL returns the domain:port scheme. Port is omitted if 80.
func (c Config) URL() (u *url.URL) {
	u = &url.URL{}
	if c.HTTPS {
		u.Scheme = "https"
	} else {
		u.Scheme = "http"
	}
	// Fallback to the non proxy domain and ports
	domain := c.ProxyDomain
	if domain == "" {
		domain = c.Domain
	}
	port := c.ProxyPort
	if port == 0 {
		port = c.Port
	}
	if port == 80 {
		u.Host = domain
	} else {
		u.Host = fmt.Sprintf("%s:%d", domain, port)
	}
	return
}

// FullAddress returns the scheme, domain, port, and host
func (c Config) FullAddress() string {
	return c.URL().String()
}

// Parse will create a Config using the file settings.json in the
// current directory.
func Parse() (Config, error) {
	return ParseFile("./settings.json")
}

// ParseFile will create a Config using the file at the given path.
func ParseFile(filename string) (Config, error) {
	f, err := os.Open(filename)
	if err != nil {
		return Config{}, err
	}
	return parse(f)
}

// TODO What about default values other than the cookie? Leave to user?
func parse(f io.Reader) (Config, error) {
	c := Config{
		Cookie: DefaultCookie,
	}
	contents, err := ioutil.ReadAll(f)
	if err != nil {
		return c, err
	}
	if err = json.Unmarshal(contents, &c); err != nil {
		return c, err
	}
	return c, nil
}

// Default is a basic configuration with insecure values. It will return the
// Address localhost:8080
var Default = Config{
	Cookie:    DefaultCookie,
	Port:      8080,
	StaticURL: "/static/",
	Metadata:  Metadata{},
}

// DefaultConfig will return a basic configuration with insecure values. It
// allows the specification of a secret key.
func DefaultConfig(key string) Config {
	return Config{
		Cookie:    DefaultCookie,
		Port:      8080,
		SecretKey: key,
		StaticURL: "/static/",
	}
}
