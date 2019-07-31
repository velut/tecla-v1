package client

import (
	"errors"
	"net/http"
	"sync"

	"github.com/rakyll/statik/fs"
)

var (
	once   sync.Once
	assets http.FileSystem
)

// Assets returns the filesystem containing the client's assets.
func Assets() (http.FileSystem, error) {
	once.Do(func() {
		assets, _ = fs.New()
	})

	noAssets := assets == nil
	if noAssets {
		return nil, errors.New("client's assets not found")
	}

	return assets, nil
}
