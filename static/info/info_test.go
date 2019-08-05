package info

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAppInfo(t *testing.T) {
	assert := assert.New(t)
	testName := "TestAppInfo"

	got := AppInfo()

	want := &Info{
		Name:        name,
		Description: description,
		Homepage:    homepage,
		Repository:  repository,
		Version:     version,
		Commit:      commit,
		Copyright:   copyright,
		Warranty:    warranty,
		License:     license,
	}

	assert.Equal(want, got, testName)
	assert.NotEqual(noVersion, got.Version, testName)
	assert.NotEqual(noCommit, got.Commit, testName)
	assert.NotEqual(noLicense, got.License, testName)
	assert.NotEmpty(got.License, testName)
}
