package info

import (
	"sync"
)

// Placeholders
const (
	noVersion = "version not set"
	noCommit  = "commit not set"
	noLicense = "license not set"
)

// App info
var (
	name        = "Tecla"
	description = "The interactive file organizer"
	homepage    = "https://github.com/velut/tecla"
	repository  = "https://github.com/velut/tecla"
	version     = noVersion // Set in gitInfo.go
	commit      = noCommit  // Set in gitInfo.go
	copyright   = "Copyright (c) 2019 Edoardo Scibona"
	warranty    = "This program comes with ABSOLUTELY NO WARRANTY"
	license     = noLicense // Read from static filesystem
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
	Warranty    string `json:"warranty"`
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
		Warranty:    warranty,
		License:     license,
	}
}

func setLicense() {
	license = _escFSMustString(false, "/LICENSE")
}
