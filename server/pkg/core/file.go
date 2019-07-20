package core

// File represents a regular file managed by the organizer.
type File struct {
	ID   int64  `json:"id"`
	Name string `json:"name"`
	Dir  string `json:"dir"`
	Path string `json:"path"`
	Ext  string `json:"ext"`
	Size int64  `json:"size"`
	URL  string `json:"url"`
}
