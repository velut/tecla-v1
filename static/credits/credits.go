package credits

import (
	"encoding/json"
	"fmt"
	"net/http"
	"sort"
	"strings"
	"sync"

	"github.com/rakyll/statik/fs"
)

// Credits contains information about the software included in the application.
type Credits struct {
	Credits []*Library `json:"credits"`
}

// Library represents an external software library.
type Library struct {
	Name     string `json:"name"`
	Homepage string `json:"homepage"`
	License  string `json:"license"`
}

var (
	once      sync.Once
	libraries []*Library
)

// AppCredits returns the credits for the application.
func AppCredits() *Credits {
	once.Do(setLibraries)

	return &Credits{
		Credits: libraries,
	}
}

func setLibraries() {
	assets, err := fs.New()
	if err != nil {
		return
	}

	clientLibraries, err := readClientLibraries(assets)
	if err != nil {
		return
	}

	serverLibraries, err := readServerLibraries(assets)
	if err != nil {
		return
	}

	allLibraries := append(clientLibraries, serverLibraries...)
	sort.Slice(allLibraries, func(i, j int) bool {
		return libDir(allLibraries[i].Name) < libDir(allLibraries[j].Name)
	})

	libraries = allLibraries
}

func readClientLibraries(assets http.FileSystem) ([]*Library, error) {
	return readLibraries(assets, "client")
}

func readServerLibraries(assets http.FileSystem) ([]*Library, error) {
	return readLibraries(assets, "server")
}

func readLibraries(assets http.FileSystem, dir string) ([]*Library, error) {
	creditsData, err := fs.ReadFile(assets, fmt.Sprintf("/%s/credits.json", dir))
	if err != nil {
		return nil, err
	}

	credits := &Credits{}
	err = json.Unmarshal(creditsData, credits)
	if err != nil {
		return nil, err
	}

	for _, lib := range credits.Credits {
		license, err := fs.ReadFile(assets, fmt.Sprintf("/%s/%s/LICENSE", dir, libDir(lib.Name)))
		if err != nil {
			return nil, err
		}
		lib.License = string(license)
	}

	return credits.Credits, nil
}

func libDir(libName string) string {
	// Transform "My Lib" into "my-lib"
	return strings.ReplaceAll(strings.ToLower(libName), " ", "-")
}
