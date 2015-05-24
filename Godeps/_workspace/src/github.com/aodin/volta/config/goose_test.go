package config

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParseTestYAML(t *testing.T) {
	assert := assert.New(t)

	c, err := ParseGooseDatabase("./test_fixtures/dbconf.yml", "test")
	assert.Nil(err)
	assert.Equal("postgres", c.Driver)
	assert.Equal("host=localhost port=5432 dbname=db_test user=test password=bad sslmode=disable", c.Credentials())
}
