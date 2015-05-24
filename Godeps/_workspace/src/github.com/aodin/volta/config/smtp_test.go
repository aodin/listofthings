package config

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSMTPConfig(t *testing.T) {
	assert := assert.New(t)

	assert.Equal("<test>", SMTPConfig{From: "test"}.FromAddress())
	assert.Equal(
		`"alias" <test>`,
		SMTPConfig{Alias: "alias", From: "test"}.FromAddress(),
	)

	assert.Equal("l:1234", SMTPConfig{Host: "l", Port: 1234}.Address())
}
