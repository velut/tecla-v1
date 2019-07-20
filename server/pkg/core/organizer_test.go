package core

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestNewOrganizer(t *testing.T) {
	assert := assert.New(t)

	got := NewOrganizer()
	assert.NotNil(got, "TestNewOrganizer")
}

func TestOrganizer_LoadConfig(t *testing.T) {
	assert := assert.New(t)

	dir1, err := ioutil.TempDir("", "dir")
	assert.Nil(err)
	defer os.RemoveAll(dir1)

	wd, err := os.Getwd()
	assert.Nil(err)
	testdataDir := filepath.Join(filepath.Dir(filepath.Dir(wd)), "testdata", "directory_tree")

	type args struct {
		config *Config
	}
	tests := []struct {
		name    string
		o       *Organizer
		args    args
		want    *OrganizerStatus
		wantErr bool
	}{
		{
			"empty config",
			NewOrganizer(),
			args{
				emptyConfig(),
			},
			nil,
			true,
		},
		{
			"config with empty source directory",
			NewOrganizer(),
			args{
				configWithSrcDir(dir1),
			},
			nil,
			true,
		},
		{
			"config with non-empty source directory",
			NewOrganizer(),
			args{
				configWithSrcDir(testdataDir),
			},
			&OrganizerStatus{
				Config: configWithSrcDir(testdataDir),
				CurrentFile: &File{
					ID:   int64(1),
					Name: "10.gif",
					Dir:  testdataDir,
					Path: filepath.Join(testdataDir, "10.gif"),
					Ext:  ".gif",
					Size: 799,
					URL:  defaultFileServerAddr + "/10.gif",
				},
				CurrentFileIndex: 0,
				NumFiles:         2,
			},
			false,
		},
		{
			"config with non-empty source directory and subdirectories",
			NewOrganizer(),
			args{
				configWithSrcDirAndSubDirs(testdataDir),
			},
			&OrganizerStatus{
				Config: configWithSrcDirAndSubDirs(testdataDir),
				CurrentFile: &File{
					ID:   int64(1),
					Name: "10.gif",
					Dir:  testdataDir,
					Path: filepath.Join(testdataDir, "10.gif"),
					Ext:  ".gif",
					Size: 799,
					URL:  defaultFileServerAddr + "/10.gif",
				},
				CurrentFileIndex: 0,
				NumFiles:         8,
			},
			false,
		},
	}
	for _, tt := range tests {
		got, gotErr := tt.o.LoadConfig(tt.args.config)
		assert.Equal(tt.wantErr, gotErr != nil, tt.name)
		assert.Equal(tt.want, got, tt.name)

		got, gotErr = tt.o.DropConfig()
		assert.Nil(gotErr, tt.name)
		assert.Equal(&OrganizerStatus{}, got, tt.name)
	}
}

func TestOrganizer_DropConfigWait(t *testing.T) {
	assert := assert.New(t)
	name := "TestOrganizer_DropConfigWait"

	o := NewOrganizer()

	got, gotErr := o.DropConfigWait()
	assert.Nil(gotErr, name)
	assert.Equal(&OrganizerStatus{}, got, name)
}

func TestOrganizer_DropConfig(t *testing.T) {
	assert := assert.New(t)
	name := "TestOrganizer_DropConfig"

	o := NewOrganizer()

	got, gotErr := o.DropConfig()
	assert.Nil(gotErr, name)
	assert.Equal(&OrganizerStatus{}, got, name)
}

func TestOrganizer_OrganizerStatus(t *testing.T) {
	assert := assert.New(t)
	name := "TestOrganizer_OrganizerStatus"

	o := NewOrganizer()

	got, gotErr := o.OrganizerStatus()
	assert.Nil(gotErr, name)
	assert.Equal(&OrganizerStatus{}, got, name)
}

func TestOrganizer_HandleHotkey(t *testing.T) {
	assert := assert.New(t)

	dir1, err := ioutil.TempDir("", "dir")
	assert.Nil(err)
	defer os.RemoveAll(dir1)

	dir2, err := ioutil.TempDir("", "dir")
	assert.Nil(err)
	defer os.RemoveAll(dir2)

	wd, err := os.Getwd()
	assert.Nil(err)
	testdataDir := filepath.Join(filepath.Dir(filepath.Dir(wd)), "testdata", "directory_tree")

	type args struct {
		hotkey string
	}
	tests := []struct {
		name    string
		o       *Organizer
		args    args
		want    *OrganizerStatus
		wantErr bool
	}{
		{
			"organizer without config",
			NewOrganizer(),
			args{
				"x",
			},
			&OrganizerStatus{},
			false,
		},
		{
			"organizer without config, skip file",
			NewOrganizer(),
			args{
				" ",
			},
			&OrganizerStatus{},
			false,
		},
		{
			"organizer with empty config",
			organizerWithConfig(emptyConfig()),
			args{
				"x",
			},
			&OrganizerStatus{},
			false,
		},
		{
			"organizer with empty config, skip file",
			organizerWithConfig(emptyConfig()),
			args{
				" ",
			},
			&OrganizerStatus{},
			false,
		},
		{
			"organizer with config with empty source directory",
			organizerWithConfig(configWithSrcDir(dir1)),
			args{
				"x",
			},
			&OrganizerStatus{},
			false,
		},
		{
			"organizer with config with empty source directory, skip file",
			organizerWithConfig(configWithSrcDir(dir1)),
			args{
				" ",
			},
			&OrganizerStatus{},
			false,
		},
		{
			"organizer with config with non-empty source directory, invalid hotkey",
			organizerWithConfig(configWithSrcDirAndDstDir(testdataDir, dir2)),
			args{
				"a",
			},
			&OrganizerStatus{
				Config: configWithSrcDirAndDstDir(testdataDir, dir2),
				CurrentFile: &File{
					ID:   int64(1),
					Name: "10.gif",
					Dir:  testdataDir,
					Path: filepath.Join(testdataDir, "10.gif"),
					Ext:  ".gif",
					Size: 799,
					URL:  defaultFileServerAddr + "/10.gif",
				},
				CurrentFileIndex: 0,
				NumFiles:         2,
			},
			false,
		},
		{
			"organizer with config with non-empty source directory, valid hotkey",
			organizerWithConfig(configWithSrcDirAndDstDir(testdataDir, dir2)),
			args{
				"x",
			},
			&OrganizerStatus{
				Config: configWithSrcDirAndDstDir(testdataDir, dir2),
				CurrentFile: &File{
					ID:   int64(2),
					Name: "20.gif",
					Dir:  testdataDir,
					Path: filepath.Join(testdataDir, "20.gif"),
					Ext:  ".gif",
					Size: 799,
					URL:  defaultFileServerAddr + "/20.gif",
				},
				CurrentFileIndex: 1,
				NumFiles:         2,
			},
			false,
		},
		{
			"organizer with config with non-empty source directory, skip file",
			organizerWithConfig(configWithSrcDirAndDstDir(testdataDir, dir2)),
			args{
				" ",
			},
			&OrganizerStatus{
				Config: configWithSrcDirAndDstDir(testdataDir, dir2),
				CurrentFile: &File{
					ID:   int64(2),
					Name: "20.gif",
					Dir:  testdataDir,
					Path: filepath.Join(testdataDir, "20.gif"),
					Ext:  ".gif",
					Size: 799,
					URL:  defaultFileServerAddr + "/20.gif",
				},
				CurrentFileIndex: 1,
				NumFiles:         2,
			},
			false,
		},
	}
	for _, tt := range tests {
		got, gotErr := tt.o.HandleHotkey(tt.args.hotkey)
		assert.Equal(tt.wantErr, gotErr != nil, tt.name)
		assert.Equal(tt.want, got, tt.name)
	}
}

func TestOrganizerInteractionCopy(t *testing.T) {
	assert := assert.New(t)
	name := "TestOrganizerInteractionCopy"

	dir1, err := ioutil.TempDir("", "dir")
	assert.Nil(err)
	defer os.RemoveAll(dir1)

	wd, err := os.Getwd()
	assert.Nil(err)
	testdataDir := filepath.Join(filepath.Dir(filepath.Dir(wd)), "testdata", "directory_tree")

	o := NewOrganizer()
	config := configWithSrcDirAndDstDir(testdataDir, dir1)

	// Load config
	status, err := o.LoadConfig(config)
	wantStatus := &OrganizerStatus{
		Config: config,
		CurrentFile: &File{
			ID:   int64(1),
			Name: "10.gif",
			Dir:  testdataDir,
			Path: filepath.Join(testdataDir, "10.gif"),
			Ext:  ".gif",
			Size: 799,
			URL:  defaultFileServerAddr + "/10.gif",
		},
		CurrentFileIndex: 0,
		NumFiles:         2,
	}
	assert.Equal(wantStatus, status, name)
	assert.Nil(err, name)

	// Ignore invalid hotkey
	status, err = o.HandleHotkey("a")
	assert.Equal(wantStatus, status, name)
	assert.Nil(err, name)

	// Skip first file
	status, err = o.HandleHotkey(" ")
	wantStatus.CurrentFile = &File{
		ID:   int64(2),
		Name: "20.gif",
		Dir:  testdataDir,
		Path: filepath.Join(testdataDir, "20.gif"),
		Ext:  ".gif",
		Size: 799,
		URL:  defaultFileServerAddr + "/20.gif",
	}
	wantStatus.CurrentFileIndex = 1
	assert.Equal(wantStatus, status, name)
	assert.Nil(err, name)

	// Ignore invalid hotkey
	status, err = o.HandleHotkey("a")
	assert.Equal(wantStatus, status, name)
	assert.Nil(err, name)

	// Operate on second file
	status, err = o.HandleHotkey("x")
	wantStatus.CurrentFile = nil
	wantStatus.CurrentFileIndex = 2
	assert.Equal(wantStatus, status, name)
	assert.Nil(err, name)
	time.Sleep(300 * time.Millisecond)
	assert.FileExists(filepath.Join(dir1, "20.gif"))

	// No more files, ignore valid hotkey
	status, err = o.HandleHotkey("x")
	assert.Equal(wantStatus, status, name)
	assert.Nil(err, name)

	// No more files, ignore skip hotkey
	status, err = o.HandleHotkey(" ")
	assert.Equal(wantStatus, status, name)
	assert.Nil(err, name)

	// No more files, ignore invalid hotkey
	status, err = o.HandleHotkey("a")
	assert.Equal(wantStatus, status, name)
	assert.Nil(err, name)

	// Drop config and wait for operations to finish
	status, err = o.DropConfigWait()
	wantStatus = &OrganizerStatus{}
	assert.Equal(wantStatus, status, name)
	assert.Nil(err, name)
}

func TestOrganizerInteractionMove(t *testing.T) {
	assert := assert.New(t)
	name := "TestOrganizerInteractionMove"

	dir1, err := ioutil.TempDir("", "dir")
	assert.Nil(err)
	file1Path := filepath.Join(dir1, "file1.txt")
	err = ioutil.WriteFile(file1Path, []byte("123"), 0644)
	assert.Nil(err)
	assert.FileExists(file1Path)
	file2Path := filepath.Join(dir1, "file2.txt")
	err = ioutil.WriteFile(file2Path, []byte("123"), 0644)
	assert.Nil(err)
	assert.FileExists(file2Path)
	defer os.RemoveAll(dir1)

	dir2, err := ioutil.TempDir("", "dir")
	assert.Nil(err)
	defer os.RemoveAll(dir2)

	o := NewOrganizer()
	config := configWithSrcDirAndDstDirMove(dir1, dir2)

	// Load config
	status, err := o.LoadConfig(config)
	wantStatus := &OrganizerStatus{
		Config: config,
		CurrentFile: &File{
			ID:   int64(1),
			Name: "file1.txt",
			Dir:  dir1,
			Path: filepath.Join(dir1, "file1.txt"),
			Ext:  ".txt",
			Size: 3,
			URL:  defaultFileServerAddr + "/file1.txt",
		},
		CurrentFileIndex: 0,
		NumFiles:         2,
	}
	assert.Equal(wantStatus, status, name)
	assert.Nil(err, name)

	// Ignore invalid hotkey
	status, err = o.HandleHotkey("a")
	assert.Equal(wantStatus, status, name)
	assert.Nil(err, name)

	// Skip first file
	status, err = o.HandleHotkey(" ")
	wantStatus.CurrentFile = &File{
		ID:   int64(2),
		Name: "file2.txt",
		Dir:  dir1,
		Path: filepath.Join(dir1, "file2.txt"),
		Ext:  ".txt",
		Size: 3,
		URL:  defaultFileServerAddr + "/file2.txt",
	}
	wantStatus.CurrentFileIndex = 1
	assert.Equal(wantStatus, status, name)
	assert.Nil(err, name)

	// Ignore invalid hotkey
	status, err = o.HandleHotkey("a")
	assert.Equal(wantStatus, status, name)
	assert.Nil(err, name)

	// Operate on second file
	status, err = o.HandleHotkey("x")
	wantStatus.CurrentFile = nil
	wantStatus.CurrentFileIndex = 2
	assert.Equal(wantStatus, status, name)
	assert.Nil(err, name)
	time.Sleep(300 * time.Millisecond)
	assert.FileExists(filepath.Join(dir2, "file2.txt"))

	// No more files, ignore valid hotkey
	status, err = o.HandleHotkey("x")
	assert.Equal(wantStatus, status, name)
	assert.Nil(err, name)

	// No more files, ignore skip hotkey
	status, err = o.HandleHotkey(" ")
	assert.Equal(wantStatus, status, name)
	assert.Nil(err, name)

	// No more files, ignore invalid hotkey
	status, err = o.HandleHotkey("a")
	assert.Equal(wantStatus, status, name)
	assert.Nil(err, name)

	// Drop config, cancel pending operations
	status, err = o.DropConfig()
	wantStatus = &OrganizerStatus{}
	assert.Equal(wantStatus, status, name)
	assert.Nil(err, name)
}

func organizerWithConfig(config *Config) *Organizer {
	o := NewOrganizer()
	_, _ = o.LoadConfig(config)
	return o
}

func emptyConfig() *Config {
	return &Config{
		ID:   0,
		Name: "",
		Src: &ConfigSrc{
			Dir:            "",
			IncludeSubdirs: false,
			DefaultOpType:  OpTypeCopy,
		},
		Dst: &ConfigDst{
			Dirs: []DstDir{},
		},
		Ops: &ConfigOps{
			NumWorkers: 0,
			MaxTries:   0,
		},
	}
}

func configWithSrcDir(dir string) *Config {
	return &Config{
		ID:   0,
		Name: "",
		Src: &ConfigSrc{
			Dir:            dir,
			IncludeSubdirs: false,
			DefaultOpType:  OpTypeCopy,
		},
		Dst: &ConfigDst{
			Dirs: []DstDir{},
		},
		Ops: &ConfigOps{
			NumWorkers: 0,
			MaxTries:   0,
		},
	}
}

func configWithSrcDirAndDstDir(srcDir, dstDir string) *Config {
	return &Config{
		ID:   0,
		Name: "",
		Src: &ConfigSrc{
			Dir:            srcDir,
			IncludeSubdirs: false,
			DefaultOpType:  OpTypeCopy,
		},
		Dst: &ConfigDst{
			Dirs: []DstDir{
				{
					Hotkey: "x",
					Dir:    dstDir,
				},
			},
		},
		Ops: &ConfigOps{
			NumWorkers: 1,
			MaxTries:   1,
		},
	}
}

func configWithSrcDirAndDstDirMove(srcDir, dstDir string) *Config {
	return &Config{
		ID:   0,
		Name: "",
		Src: &ConfigSrc{
			Dir:            srcDir,
			IncludeSubdirs: false,
			DefaultOpType:  OpTypeMove,
		},
		Dst: &ConfigDst{
			Dirs: []DstDir{
				{
					Hotkey: "x",
					Dir:    dstDir,
				},
			},
		},
		Ops: &ConfigOps{
			NumWorkers: 1,
			MaxTries:   1,
		},
	}
}

func configWithSrcDirAndSubDirs(dir string) *Config {
	return &Config{
		ID:   0,
		Name: "",
		Src: &ConfigSrc{
			Dir:            dir,
			IncludeSubdirs: true,
			DefaultOpType:  OpTypeCopy,
		},
		Dst: &ConfigDst{
			Dirs: []DstDir{},
		},
		Ops: &ConfigOps{
			NumWorkers: 0,
			MaxTries:   0,
		},
	}
}
