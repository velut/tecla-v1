package credits

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

const numCredits = 21

func TestAppCredits(t *testing.T) {
	assert := assert.New(t)
	testName := "TestAppCredits"

	got := AppCredits()
	assert.NotNil(got.Credits, testName)
	assert.Len(got.Credits, numCredits, testName)
	for i := 0; i < numCredits; i++ {
		assert.NotEmpty(got.Credits[i].License, testName)
	}
}
