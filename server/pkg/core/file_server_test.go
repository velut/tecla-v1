package core

import (
	"io/ioutil"
	"net/http"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewFileServer(t *testing.T) {
	assert := assert.New(t)

	dir, err := ioutil.TempDir("", "dir")
	assert.Nil(err)
	defer os.RemoveAll(dir)

	fs := NewFileServer(dir)
	assert.NotNil(fs)
	defer fs.Close()

	resp, err := http.Get(defaultFileServerAddr)
	assert.Nil(err)
	resp.Body.Close()
}

func TestFileServer_Close(t *testing.T) {
	assert := assert.New(t)

	dir, err := ioutil.TempDir("", "dir")
	assert.Nil(err)
	defer os.RemoveAll(dir)

	fs := NewFileServer(dir)
	assert.NotNil(fs)

	resp, err := http.Get(defaultFileServerAddr)
	assert.Nil(err)
	resp.Body.Close()

	err = fs.Close()
	assert.Nil(err)
}
