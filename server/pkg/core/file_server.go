package core

import (
	"net/http"
	"time"
)

const defaultFileServerAddr = "http://localhost:5921"
const defaultFileServerPort = ":5921"

// FileServer represents an HTTP server that serves a static directory.
type FileServer struct {
	server *http.Server
}

// NewFileServer creates a new FileServer serving the files present in the given directory.
func NewFileServer(dir string) *FileServer {
	fs := &FileServer{
		server: &http.Server{
			Addr:    defaultFileServerPort,
			Handler: http.FileServer(http.Dir(dir)),
		},
	}
	go func() {
		_ = fs.server.ListenAndServe()
	}()
	time.Sleep(100 * time.Millisecond)
	return fs
}

// Close stops the FileServer.
func (s *FileServer) Close() error {
	return s.server.Close()
}
