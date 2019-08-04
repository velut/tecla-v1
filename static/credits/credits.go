package credits

import (
	"encoding/json"
	"fmt"
	"sort"
	"strings"
	"sync"
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
	clientLibraries, err := readClientLibraries()
	if err != nil {
		return
	}

	serverLibraries, err := readServerLibraries()
	if err != nil {
		return
	}

	allLibraries := append(clientLibraries, serverLibraries...)
	sort.Slice(allLibraries, func(i, j int) bool {
		return libDir(allLibraries[i].Name) < libDir(allLibraries[j].Name)
	})

	libraries = allLibraries
}

func readClientLibraries() ([]*Library, error) {
	return readLibraries("client")
}

func readServerLibraries() ([]*Library, error) {
	return readLibraries("server")
}

func readLibraries(dir string) ([]*Library, error) {
	creditsData := _escFSMustByte(false, fmt.Sprintf("/%s/credits.json", dir))

	credits := &Credits{}
	err := json.Unmarshal(creditsData, credits)
	if err != nil {
		return nil, err
	}

	for _, lib := range credits.Credits {
		lib.License = _escFSMustString(false, fmt.Sprintf("/%v/%v/LICENSE", dir, libDir(lib.Name)))
	}

	return credits.Credits, nil
}

func libDir(libName string) string {
	// Transform "My Lib" into "my-lib"
	return strings.ReplaceAll(strings.ToLower(libName), " ", "-")
}
