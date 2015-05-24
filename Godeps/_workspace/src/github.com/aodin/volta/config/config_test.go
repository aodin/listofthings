package config

import (
	"net/url"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestConfig(t *testing.T) {
	assert := assert.New(t)

	c, err := ParseFile("./test_fixtures/settings.json")
	assert.Nil(err, "Parse file should not error")

	// Test the parent config methods
	assert.Equal("localhost:9001", c.Address())
	assert.Equal("http://localhost:9001", c.URL().String())
	assert.Equal(c.URL().String(), c.FullAddress())
	assert.Equal("1.0.0", c.Metadata.Get("version"))

	c.ProxyDomain = "example.com"
	c.ProxyPort = 3000
	assert.Equal("http://example.com:3000", c.FullAddress())
	assert.Equal(
		&url.URL{Scheme: "http", Host: "example.com:3000"},
		c.URL(),
	)

	c.ProxyPort = 80
	c.HTTPS = true
	assert.Equal("https://example.com", c.FullAddress())
	assert.Equal(
		&url.URL{Scheme: "https", Host: "example.com"},
		c.URL(),
	)

	// Test the SMTP config methods
	assert.Equal(`"Example User" <no_reply@example.com>`, c.SMTP.FromAddress())
	assert.Equal("example.com:587", c.SMTP.Address())

	// Test the database config methods when values are missing
	assert.Equal(
		"host=localhost port=5432 dbname=db user=pg password=pass",
		c.Database.Credentials(),
	)

	// Test the default cookie settings
	assert.Equal(336*time.Hour, c.Cookie.Age)
	assert.Equal("", c.Cookie.Domain)
	assert.Equal(false, c.Cookie.HttpOnly)
	assert.Equal("sessionid", c.Cookie.Name)
	assert.Equal("/", c.Cookie.Path)
	assert.Equal(false, c.Cookie.Secure)

	// TODO Test custom cookie settings
}
