package info

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAppInfo(t *testing.T) {
	assert := assert.New(t)

	got := AppInfo()

	want := &Info{
		Name:        name,
		Description: description,
		Homepage:    homepage,
		Repository:  repository,
		Version:     "not set",
		Commit:      "not set",
		Copyright:   copyright,
		NoWarranty:  noWarranty,
		License:     "not set",
	}

	assert.Equal(want, got, "TestAppInfo")
}
