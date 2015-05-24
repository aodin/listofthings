package config

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMetadata(t *testing.T) {
	assert := assert.New(t)

	m := Metadata{"version": "1.0.0"}
	assert.Equal([]string{"version"}, m.Keys())
	assert.Equal([]string{"1.0.0"}, m.Values())
	assert.True(m.Has("version"))
	assert.Equal("1.0.0", m.Get("version"))
}
