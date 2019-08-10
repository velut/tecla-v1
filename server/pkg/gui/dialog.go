package gui

import (
	"errors"
	"sync/atomic"

	"github.com/gen2brain/dlgs"
)

const windowName = "Tecla - The interactive file organizer"

var (
	selectDirectoryLock uint32
)

// SelectDirectory opens a directory selection dialog
// and returns the path of the selected directory.
// Only one dialog can be open at a time.
func SelectDirectory() (string, error) {
	if !atomic.CompareAndSwapUint32(&selectDirectoryLock, 0, 1) {
		return "", errors.New("directory selection dialog is already open")
	}
	defer atomic.StoreUint32(&selectDirectoryLock, 0)

	dir, ok, err := dlgs.File("Select a directory", "", windowName, true)
	if err != nil {
		return "", err
	}
	if !ok {
		return "", errors.New("could not select directory")
	}

	return dir, nil
}
