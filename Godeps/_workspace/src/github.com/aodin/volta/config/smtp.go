package config

import "fmt"

// SMTPConfig contains the fields needed to connect to a SMTP server.
type SMTPConfig struct {
	Port     int64  `json:"port"`
	User     string `json:"user"`
	Password string `json:"password"`
	Host     string `json:"host"`
	From     string `json:"from"`
	Alias    string `json:"alias"`
}

// FromAddress creates a string suitable for use in an Email's From header.
func (c SMTPConfig) FromAddress() string {
	if c.Alias != "" {
		return fmt.Sprintf(`"%s" <%s>`, c.Alias, c.From)
	}
	return fmt.Sprintf("<%s>", c.From)
}

// Address will return a string of the host and port separated by a colon.
func (c SMTPConfig) Address() string {
	return fmt.Sprintf("%s:%d", c.Host, c.Port)
}
