package info

import (
	"sync"

	"github.com/rakyll/statik/fs"
)

// App info
var (
	name        = "Tecla"
	description = "The interactive file organizer"
	homepage    = "https://github.com/velut/tecla"
	repository  = "https://github.com/velut/tecla"
	version     = "version not set" // Set by ldflags
	commit      = "commit not set"  // Set by ldflags
	copyright   = "Copyright (c) 2019 Edoardo Scibona"
	noWarranty  = "This program comes with ABSOLUTELY NO WARRANTY"
	license     = "license not set" // Read from static filesystem
)

// Info lists information about the application.
type Info struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	Homepage    string `json:"homepage"`
	Repository  string `json:"repository"`
	Version     string `json:"version"`
	Commit      string `json:"commit"`
	Copyright   string `json:"copyright"`
	NoWarranty  string `json:"noWarranty"`
	License     string `json:"license"`
}

var once sync.Once

// AppInfo returns information about the application.
func AppInfo() *Info {
	once.Do(setLicense)

	return &Info{
		Name:        name,
		Description: description,
		Homepage:    homepage,
		Repository:  repository,
		Version:     version,
		Commit:      commit,
		Copyright:   copyright,
		NoWarranty:  noWarranty,
		License:     license,
	}
}

func setLicense() {
	assets, err := fs.New()
	if err != nil {
		return
	}

	l, err := fs.ReadFile(assets, "/LICENSE")
	if err != nil {
		return
	}

	license = string(l)
}
