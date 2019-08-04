package info

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAppInfo(t *testing.T) {
	assert := assert.New(t)
	testName := "TestAppInfo"

	noVersion := version
	noCommit := commit
	noLicense := license

	got := AppInfo()

	want := &Info{
		Name:        name,
		Description: description,
		Homepage:    homepage,
		Repository:  repository,
		Version:     noVersion, // Set by ldflags
		Commit:      noCommit,  // Set by ldflags
		Copyright:   copyright,
		NoWarranty:  noWarranty,
		License:     license,
	}

	assert.Equal(want, got, testName)
	assert.NotEqual(noLicense, got.License, testName)
	assert.NotEmpty(got.License, testName)
}
