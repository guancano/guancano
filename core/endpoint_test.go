package core

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestParse(t *testing.T) {
	prefix, url, optStr, _ := Parse("mock:test1")
	assert.Equal(t, "mock", prefix)
	assert.Equal(t, "test1", url)

	prefix, url, _, _ = Parse("mock:")
	assert.Equal(t, "mock", prefix)
	assert.Equal(t, "", url)

	prefix, url, _, _ = Parse(":")
	assert.Equal(t, "", prefix)
	assert.Equal(t, "", url)

	prefix, url, _, _ = Parse("")
	assert.Equal(t, "", prefix)
	assert.Equal(t, "", url)

	prefix, url, optStr, _ = Parse("")
	assert.Equal(t, "", prefix)
	assert.Equal(t, "", url)
	assert.Equal(t, "", optStr)

	prefix, url, optStr, _ = Parse("?dogfly=12")
	assert.Equal(t, "", prefix)
	assert.Equal(t, "", url)
	assert.Equal(t, "", optStr)

	prefix, url, optStr, _ = Parse(":?")
	assert.Equal(t, "", prefix)
	assert.Equal(t, "", url)
	assert.Equal(t, "", optStr)

	prefix, url, optStr, _ = Parse("prefix:/some/path/url?option1=option&option2=option")
	assert.Equal(t, "prefix", prefix)
	assert.Equal(t, "/some/path/url", url)
	assert.Equal(t, "option1=option&option2=option", optStr)
}
