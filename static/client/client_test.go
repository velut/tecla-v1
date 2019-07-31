package client

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAssets(t *testing.T) {
	assert := assert.New(t)
	testName := "TestAssets"

	got, gotErr := Assets()

	assert.NoError(gotErr, testName)
	assert.NotNil(got, testName)
}
