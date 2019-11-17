package core

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"sync"
	"time"

	"github.com/gammazero/workerpool"
	"github.com/velut/fsutils-go/fs"
)

// Organizer represents the organizer that handles files and file operations.
type Organizer struct {
	mutex     sync.Mutex
	organizer *organizer
}

type organizer struct {
	config           *Config
	files            Files
	currentFileIndex int
	fileServer       *FileServer
	workerPool       *workerpool.WorkerPool
}

// Files represents a collection of files.
type Files []*File

// OrganizerStatus represents the status of the organizer.
type OrganizerStatus struct {
	Config           *Config `json:"config"`
	CurrentFile      *File   `json:"currentFile"`
	CurrentFileIndex int     `json:"currentFileIndex"`
	NumFiles         int     `json:"numFiles"`
}

// NewOrganizer creates a new Organizer.
func NewOrganizer() *Organizer {
	return &Organizer{
		organizer: newOrganizer(),
	}
}

func newOrganizer() *organizer {
	return &organizer{}
}

// RestoreConfig TODO:
func (o *Organizer) RestoreConfig() (*Config, error) {
	o.mutex.Lock()
	defer o.mutex.Unlock()

	// Get latest config
	ucDir, err := os.UserConfigDir()
	if err != nil {
		return nil, err
	}

	teclaDir := filepath.Join(ucDir, "tecla")
	latestFile := filepath.Join(teclaDir, "tecla-latest-config.json")
	configJSON, err := ioutil.ReadFile(latestFile)
	if err != nil {
		return nil, err
	}
	config := &Config{}
	err = json.Unmarshal(configJSON, config)
	if err != nil {
		return nil, err
	}

	return config, nil
}

// LoadConfig loads the given configuration, which must be valid, starting the organizer.
func (o *Organizer) LoadConfig(config *Config) (*OrganizerStatus, error) {
	o.mutex.Lock()
	defer o.mutex.Unlock()

	o.dropConfig()
	if err := o.loadConfig(config); err != nil {
		o.dropConfig()
		return nil, err
	}

	// Save as latest config
	configJSON, _ := json.Marshal(config)
	ucDir, _ := os.UserConfigDir()
	teclaDir := filepath.Join(ucDir, "tecla")
	_ = os.MkdirAll(teclaDir, 0700)
	latestFile := filepath.Join(teclaDir, "tecla-latest-config.json")
	_ = ioutil.WriteFile(latestFile, configJSON, 0644)
	//

	return o.organizerStatus()
}

func (o *Organizer) loadConfig(config *Config) error {
	return o.organizer.loadConfig(config)
}

func (o *organizer) loadConfig(config *Config) error {
	o.config = config

	if err := o.gatherFiles(); err != nil {
		return err
	}

	o.startFileServer()
	o.startWorkerPool()

	return nil
}

func (o *organizer) gatherFiles() error {
	configSrc := o.config.Src
	srcDir := configSrc.Dir
	includeSubdirs := configSrc.IncludeSubdirs

	fileInfos, err := fs.ReadDir(srcDir, &fs.ReadDirOptions{
		IncludeSubdirs: includeSubdirs,
	})
	if err != nil {
		return err
	}
	noFiles := len(fileInfos) == 0
	if noFiles {
		return errors.New("no files to organize")
	}

	o.files = make(Files, len(fileInfos))
	for i, fi := range fileInfos {
		o.files[i] = &File{
			ID:   int64(i + 1),
			Name: fi.Name,
			Dir:  fi.Dir,
			Path: fi.Path,
			Ext:  fi.Ext,
			Size: fi.Size,
			URL:  fileURL(srcDir, fi.Path),
		}
	}

	return nil
}

func fileURL(srcDir, filePath string) string {
	rel, _ := filepath.Rel(srcDir, filePath)
	serverPath := filepath.ToSlash(rel)
	return defaultFileServerAddr + "/" + serverPath
}

func (o *organizer) startFileServer() {
	o.fileServer = NewFileServer(o.config.Src.Dir)
}

func (o *organizer) startWorkerPool() {
	o.workerPool = workerpool.New(o.config.Ops.NumWorkers)
}

// DropConfigWait removes the current configuration, if any, stopping the organizer.
// All submitted operations, pending or in progress, are completed.
func (o *Organizer) DropConfigWait() (*OrganizerStatus, error) {
	o.mutex.Lock()
	defer o.mutex.Unlock()

	o.dropConfigWait()

	return o.organizerStatus()
}

func (o *Organizer) dropConfigWait() {
	o.organizer.stopWait()
	o.organizer = newOrganizer()
}

func (o *organizer) stopWait() {
	o.stopFileServer()
	o.stopWorkerPoolWait()
}

func (o *organizer) stopWorkerPoolWait() {
	if o.workerPool != nil {
		o.workerPool.StopWait()
	}
}

// DropConfig removes the current configuration, if any, stopping the organizer.
// In progress operations are completed, pending operations are discarded.
func (o *Organizer) DropConfig() (*OrganizerStatus, error) {
	o.mutex.Lock()
	defer o.mutex.Unlock()

	o.dropConfig()

	return o.organizerStatus()
}

func (o *Organizer) dropConfig() {
	o.organizer.stop()
	o.organizer = newOrganizer()
}

func (o *organizer) stop() {
	o.stopFileServer()
	o.stopWorkerPool()
}

func (o *organizer) stopFileServer() {
	if o.fileServer != nil {
		_ = o.fileServer.Close()
	}
}

func (o *organizer) stopWorkerPool() {
	if o.workerPool != nil {
		o.workerPool.Stop()
	}
}

// OrganizerStatus returns the organizer's status.
func (o *Organizer) OrganizerStatus() (*OrganizerStatus, error) {
	o.mutex.Lock()
	defer o.mutex.Unlock()
	return o.organizerStatus()
}

func (o *Organizer) organizerStatus() (*OrganizerStatus, error) {
	return o.organizer.status()
}

func (o *organizer) status() (*OrganizerStatus, error) {
	status := &OrganizerStatus{}
	if o.hasConfig() {
		status.Config = o.config
		status.CurrentFile = o.currentFile()
		status.CurrentFileIndex = o.currentFileIndex
		status.NumFiles = len(o.files)
	}
	return status, nil
}

// HandleHotkey handles the action corresponding to the pressed hotkey.
// If the hotkey is the space character, the organizer advances to the next file.
// If the hotkey is associated to a destination directory, the organizer creates
// a file operation for the current file and then advances to the next file.
// If the hotkey is unrecognized, the organizer does nothing.
func (o *Organizer) HandleHotkey(hotkey string) (*OrganizerStatus, error) {
	o.mutex.Lock()
	defer o.mutex.Unlock()

	o.handleHotkey(hotkey)

	return o.organizerStatus()
}

func (o *Organizer) handleHotkey(hotkey string) {
	o.organizer.handleHotkey(hotkey)
}

func (o *organizer) handleHotkey(hotkey string) {
	noConfig := !o.hasConfig()
	if noConfig {
		return
	}

	skipFile := hotkey == " "
	if skipFile {
		o.incrementCurrentFileIndex()
		return
	}

	if err := o.submitOperation(hotkey); err != nil {
		return
	}

	o.incrementCurrentFileIndex()
}

func (o *organizer) submitOperation(hotkey string) error {
	op, err := o.createOperation(hotkey)
	if err != nil {
		return err
	}

	o.workerPool.Submit(func() {
		// Slow down workers. If they are too fast,
		// they may try to access files still displayed in the gui;
		// this causes problems when trying to remove files
		// currently being served by the fileserver.
		time.Sleep(250 * time.Millisecond)
		o.executeOperation(op)
	})

	return nil
}

func (o *organizer) createOperation(hotkey string) (*Operation, error) {
	dstDir, ok := o.getDstDir(hotkey)
	if !ok {
		return nil, fmt.Errorf("hotkey %q not found", hotkey)
	}

	file := o.currentFile()
	if file == nil {
		return nil, fmt.Errorf("no current file found")
	}

	op := &Operation{
		ID:       file.ID,
		Op:       o.config.Src.DefaultOpType,
		SrcPath:  file.Path,
		DstPath:  filepath.Join(dstDir, file.Name),
		MaxTries: o.config.Ops.MaxTries,
	}

	return op, nil
}

func (o *organizer) executeOperation(op *Operation) {
	opType := op.Op
	srcPath := op.SrcPath
	dstPath := op.DstPath
	maxTries := op.MaxTries

	switch opType {
	case OpTypeCopy:
		_, _ = fs.CopyFileSafe(srcPath, dstPath, maxTries)
	case OpTypeMove:
		_, _ = fs.MoveFileSafe(srcPath, dstPath, maxTries)
	}
}

func (o *organizer) incrementCurrentFileIndex() {
	if o.hasCurrentFile() {
		o.currentFileIndex++
	}
}

func (o *organizer) getDstDir(hotkey string) (string, bool) {
	for _, d := range o.config.Dst.Dirs {
		if d.Hotkey == hotkey {
			return d.Dir, true
		}
	}
	return "", false
}

func (o *organizer) currentFile() *File {
	if o.hasCurrentFile() {
		return o.files[o.currentFileIndex]
	}
	return nil
}

func (o *organizer) hasCurrentFile() bool {
	return o.currentFileIndex < len(o.files)
}

func (o *organizer) hasConfig() bool {
	return o.config != nil
}
