package core

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestConfigValidationError_Error(t *testing.T) {
	assert := assert.New(t)
	tests := []struct {
		name string
		e    *ConfigValidationError
		want string
	}{
		{
			"no errors",
			&ConfigValidationError{make(errorsByKey)},
			`{"errors":{}}`,
		},
		{
			"some errors",
			&ConfigValidationError{
				errorsByKey{
					"key1": "err1",
					"key2": "err2",
					"key3": "err3",
				},
			},
			`{"errors":{"key1":"err1","key2":"err2","key3":"err3"}}`,
		},
	}
	for _, tt := range tests {
		got := tt.e.Error()
		assert.Equal(tt.want, got, tt.name)
	}
}

func TestNewConfigValidator(t *testing.T) {
	assert := assert.New(t)

	got := NewConfigValidator()
	assert.NotNil(got)
}

func TestDefaultConfigValidator_Validate(t *testing.T) {
	assert := assert.New(t)

	// Root dir
	dir1, err := ioutil.TempDir("", "dir1")
	assert.Nil(err)
	dir1File1, err := ioutil.TempFile(dir1, "dir1File1")
	assert.Nil(err)
	dir1File1.Close()
	dir1Subdir, err := ioutil.TempDir(dir1, "dir1Subdir")
	assert.Nil(err)
	defer os.RemoveAll(dir1)

	dir2, err := ioutil.TempDir("", "dir2")
	assert.Nil(err)
	defer os.RemoveAll(dir2)

	dir3, err := ioutil.TempDir("", "dir3")
	assert.Nil(err)
	defer os.RemoveAll(dir3)

	tests := []struct {
		name    string
		config  *Config
		wantErr bool
	}{
		{
			"no config",
			nil,
			true,
		},
		{
			"no config src",
			&Config{},
			true,
		},
		{
			"no config dst",
			&Config{
				Src: &ConfigSrc{},
			},
			true,
		},
		{
			"no config ops",
			&Config{
				Src: &ConfigSrc{},
				Dst: &ConfigDst{},
			},
			true,
		},
		{
			"no config name",
			&Config{
				Src: &ConfigSrc{},
				Dst: &ConfigDst{},
				Ops: &ConfigOps{},
			},
			true,
		},
		{
			"no config src dir",
			&Config{
				Name: "foo",
				Src:  &ConfigSrc{},
				Dst:  &ConfigDst{},
				Ops:  &ConfigOps{},
			},
			true,
		},
		{
			"invalid config src dir (1)",
			&Config{
				Name: "foo",
				Src: &ConfigSrc{
					Dir: "invalid",
				},
				Dst: &ConfigDst{},
				Ops: &ConfigOps{},
			},
			true,
		},
		{
			"invalid config src dir (2)",
			&Config{
				Name: "foo",
				Src: &ConfigSrc{
					Dir: dir1File1.Name(),
				},
				Dst: &ConfigDst{},
				Ops: &ConfigOps{},
			},
			true,
		},
		{
			"invalid config src dir (3)",
			&Config{
				Name: "foo",
				Src: &ConfigSrc{
					Dir: dir2,
				},
				Dst: &ConfigDst{},
				Ops: &ConfigOps{},
			},
			true,
		},
		{
			"invalid config src default op type",
			&Config{
				Name: "foo",
				Src: &ConfigSrc{
					Dir:           dir1,
					DefaultOpType: "invalid",
				},
				Dst: &ConfigDst{},
				Ops: &ConfigOps{},
			},
			true,
		},
		{
			"nil config dst dirs",
			&Config{
				Name: "foo",
				Src: &ConfigSrc{
					Dir:           dir1,
					DefaultOpType: OpTypeCopy,
				},
				Dst: &ConfigDst{},
				Ops: &ConfigOps{},
			},
			true,
		},
		{
			"empty config dst dirs",
			&Config{
				Name: "foo",
				Src: &ConfigSrc{
					Dir:           dir1,
					DefaultOpType: OpTypeCopy,
				},
				Dst: &ConfigDst{
					Dirs: []*DstDir{},
				},
				Ops: &ConfigOps{},
			},
			true,
		},
		{
			"empty config dst dirs hotkey",
			&Config{
				Name: "foo",
				Src: &ConfigSrc{
					Dir:           dir1,
					DefaultOpType: OpTypeCopy,
				},
				Dst: &ConfigDst{
					Dirs: []*DstDir{
						{"", ""},
					},
				},
				Ops: &ConfigOps{},
			},
			true,
		},
		{
			"long config dst dirs hotkey",
			&Config{
				Name: "foo",
				Src: &ConfigSrc{
					Dir:           dir1,
					DefaultOpType: OpTypeCopy,
				},
				Dst: &ConfigDst{
					Dirs: []*DstDir{
						{"tooLongHotkey", ""},
					},
				},
				Ops: &ConfigOps{},
			},
			true,
		},
		{
			"empty config dst dirs path",
			&Config{
				Name: "foo",
				Src: &ConfigSrc{
					Dir:           dir1,
					DefaultOpType: OpTypeCopy,
				},
				Dst: &ConfigDst{
					Dirs: []*DstDir{
						{"a", ""},
					},
				},
				Ops: &ConfigOps{},
			},
			true,
		},
		{
			"invalid config dst dirs path (1)",
			&Config{
				Name: "foo",
				Src: &ConfigSrc{
					Dir:           dir1,
					DefaultOpType: OpTypeCopy,
				},
				Dst: &ConfigDst{
					Dirs: []*DstDir{
						{"a", "invalid"},
					},
				},
				Ops: &ConfigOps{},
			},
			true,
		},
		{
			"invalid config dst dirs path (2)",
			&Config{
				Name: "foo",
				Src: &ConfigSrc{
					Dir:           dir1,
					DefaultOpType: OpTypeCopy,
				},
				Dst: &ConfigDst{
					Dirs: []*DstDir{
						{"a", dir1File1.Name()},
					},
				},
				Ops: &ConfigOps{},
			},
			true,
		},
		{
			"invalid config dst dirs path (3)",
			&Config{
				Name: "foo",
				Src: &ConfigSrc{
					Dir:           dir1,
					DefaultOpType: OpTypeCopy,
				},
				Dst: &ConfigDst{
					Dirs: []*DstDir{
						{"a", dir1},
					},
				},
				Ops: &ConfigOps{},
			},
			true,
		},
		{
			"invalid config dst dirs path (4)",
			&Config{
				Name: "foo",
				Src: &ConfigSrc{
					Dir:           dir1,
					DefaultOpType: OpTypeCopy,
				},
				Dst: &ConfigDst{
					Dirs: []*DstDir{
						{"a", dir1Subdir},
					},
				},
				Ops: &ConfigOps{},
			},
			true,
		},
		{
			"invalid config ops num workers (1)",
			&Config{
				Name: "foo",
				Src: &ConfigSrc{
					Dir:           dir1,
					DefaultOpType: OpTypeCopy,
				},
				Dst: &ConfigDst{
					Dirs: []*DstDir{
						{"a", dir2},
					},
				},
				Ops: &ConfigOps{
					NumWorkers: -1,
				},
			},
			true,
		},
		{
			"invalid config ops num workers (2)",
			&Config{
				Name: "foo",
				Src: &ConfigSrc{
					Dir:           dir1,
					DefaultOpType: OpTypeCopy,
				},
				Dst: &ConfigDst{
					Dirs: []*DstDir{
						{"a", dir2},
					},
				},
				Ops: &ConfigOps{
					NumWorkers: 0,
				},
			},
			true,
		},
		{
			"invalid config ops num workers (3)",
			&Config{
				Name: "foo",
				Src: &ConfigSrc{
					Dir:           dir1,
					DefaultOpType: OpTypeCopy,
				},
				Dst: &ConfigDst{
					Dirs: []*DstDir{
						{"a", dir2},
					},
				},
				Ops: &ConfigOps{
					NumWorkers: 6,
				},
			},
			true,
		},
		{
			"invalid config ops num workers (4)",
			&Config{
				Name: "foo",
				Src: &ConfigSrc{
					Dir:           dir1,
					DefaultOpType: OpTypeCopy,
				},
				Dst: &ConfigDst{
					Dirs: []*DstDir{
						{"a", dir2},
					},
				},
				Ops: &ConfigOps{
					NumWorkers: 100,
				},
			},
			true,
		},
		{
			"invalid config ops max tries (1)",
			&Config{
				Name: "foo",
				Src: &ConfigSrc{
					Dir:           dir1,
					DefaultOpType: OpTypeCopy,
				},
				Dst: &ConfigDst{
					Dirs: []*DstDir{
						{"a", dir2},
					},
				},
				Ops: &ConfigOps{
					NumWorkers: 1,
					MaxTries:   -1,
				},
			},
			true,
		},
		{
			"invalid config ops max tries (2)",
			&Config{
				Name: "foo",
				Src: &ConfigSrc{
					Dir:           dir1,
					DefaultOpType: OpTypeCopy,
				},
				Dst: &ConfigDst{
					Dirs: []*DstDir{
						{"a", dir2},
					},
				},
				Ops: &ConfigOps{
					NumWorkers: 1,
					MaxTries:   0,
				},
			},
			true,
		},
		{
			"invalid config ops max tries (3)",
			&Config{
				Name: "foo",
				Src: &ConfigSrc{
					Dir:           dir1,
					DefaultOpType: OpTypeCopy,
				},
				Dst: &ConfigDst{
					Dirs: []*DstDir{
						{"a", dir2},
					},
				},
				Ops: &ConfigOps{
					NumWorkers: 1,
					MaxTries:   1000001,
				},
			},
			true,
		},
		{
			"invalid config ops max tries (4)",
			&Config{
				Name: "foo",
				Src: &ConfigSrc{
					Dir:           dir1,
					DefaultOpType: OpTypeCopy,
				},
				Dst: &ConfigDst{
					Dirs: []*DstDir{
						{"a", dir2},
					},
				},
				Ops: &ConfigOps{
					NumWorkers: 1,
					MaxTries:   1000000000,
				},
			},
			true,
		},
		{
			"valid config (1)",
			&Config{
				Name: "foo",
				Src: &ConfigSrc{
					Dir:           dir1,
					DefaultOpType: OpTypeCopy,
				},
				Dst: &ConfigDst{
					Dirs: []*DstDir{
						{"a", dir2},
					},
				},
				Ops: &ConfigOps{
					NumWorkers: 1,
					MaxTries:   1,
				},
			},
			false,
		},
		{
			"valid config (2)",
			&Config{
				Name: "foo",
				Src: &ConfigSrc{
					Dir:           dir1,
					DefaultOpType: OpTypeMove,
				},
				Dst: &ConfigDst{
					Dirs: []*DstDir{
						{"a", dir2},
					},
				},
				Ops: &ConfigOps{
					NumWorkers: 1,
					MaxTries:   1,
				},
			},
			false,
		},
		{
			"valid config (3)",
			&Config{
				Name: "foo",
				Src: &ConfigSrc{
					Dir:           dir1,
					DefaultOpType: OpTypeMove,
				},
				Dst: &ConfigDst{
					Dirs: []*DstDir{
						{"a", dir2},
						{"b", dir2},
					},
				},
				Ops: &ConfigOps{
					NumWorkers: 1,
					MaxTries:   1,
				},
			},
			false,
		},
		{
			"valid config (4)",
			&Config{
				Name: "foo",
				Src: &ConfigSrc{
					Dir:            dir1,
					DefaultOpType:  OpTypeMove,
					IncludeSubdirs: true,
				},
				Dst: &ConfigDst{
					Dirs: []*DstDir{
						{"a", dir2},
						{"b", dir2},
					},
				},
				Ops: &ConfigOps{
					NumWorkers: 1,
					MaxTries:   1,
				},
			},
			false,
		},
		{
			"valid config (5)",
			&Config{
				Name: "foo",
				Src: &ConfigSrc{
					Dir:            dir1,
					DefaultOpType:  OpTypeCopy,
					IncludeSubdirs: true,
				},
				Dst: &ConfigDst{
					Dirs: []*DstDir{
						{"0", dir2},
						{"1", dir2},
						{"a", dir2},
						{"b", dir2},
						{"C", dir2},
						{"Д", dir3},
						{"è", dir3},
						{"β", dir3},
						{"あ", dir3},
						{"𛀀", dir3},
						{"'", dir3},
						{`\`, dir3},
					},
				},
				Ops: &ConfigOps{
					NumWorkers: 1,
					MaxTries:   1,
				},
			},
			false,
		},
		{
			"valid config (6)",
			&Config{
				Name: "foo",
				Src: &ConfigSrc{
					Dir:            dir1,
					DefaultOpType:  OpTypeCopy,
					IncludeSubdirs: true,
				},
				Dst: &ConfigDst{
					Dirs: []*DstDir{
						{"0", dir2},
						{"1", dir2},
						{"a", dir2},
						{"b", dir2},
						{"C", dir2},
						{"Д", dir3},
						{"è", dir3},
						{"β", dir3},
						{"あ", dir3},
						{"𛀀", dir3},
						{"'", dir3},
						{`\`, dir3},
					},
				},
				Ops: &ConfigOps{
					NumWorkers: 5,
					MaxTries:   1000000,
				},
			},
			false,
		},
	}
	for _, tt := range tests {
		gotErr := NewConfigValidator().ValidateConfig(tt.config)
		assert.Equal(tt.wantErr, gotErr != nil, tt.name)
	}
}

func Test_areNotSameDir(t *testing.T) {
	assert := assert.New(t)

	dir1, err := ioutil.TempDir("", "dir1")
	assert.Nil(err)
	defer os.RemoveAll(dir1)
	dir2, err := ioutil.TempDir("", "dir2")
	assert.Nil(err)
	defer os.RemoveAll(dir2)

	type args struct {
		path1 string
		path2 string
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			"invalid dirs",
			args{
				"invalid1",
				"invalid2",
			},
			false,
		},
		{
			"invalid dir 1",
			args{
				"invalid1",
				dir1,
			},
			false,
		},
		{
			"invalid dir 2",
			args{
				dir1,
				"invalid2",
			},
			false,
		},
		{
			"same dir",
			args{
				dir1,
				dir1,
			},
			false,
		},
		{
			"different dir (1)",
			args{
				dir1,
				dir2,
			},
			true,
		},
		{
			"different dir (2)",
			args{
				dir2,
				dir1,
			},
			true,
		},
	}
	for _, tt := range tests {
		got := areNotSameDir(tt.args.path1, tt.args.path2)
		assert.Equal(tt.want, got, tt.name)
	}
}

func Test_isNotChildDirOf(t *testing.T) {
	assert := assert.New(t)

	dir1, err := ioutil.TempDir("", "dir1")
	assert.Nil(err)
	defer os.RemoveAll(dir1)
	dir1Subdir, err := ioutil.TempDir(dir1, "dir1Subdir")
	assert.Nil(err)
	dir2, err := ioutil.TempDir("", "dir2")
	assert.Nil(err)
	defer os.RemoveAll(dir2)
	type args struct {
		childPath  string
		parentPath string
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			"invalid dirs",
			args{
				"invalid1",
				"invalid2",
			},
			false,
		},
		{
			"invalid dir 1",
			args{
				"invalid1",
				dir1,
			},
			false,
		},
		{
			"invalid dir 2",
			args{
				dir1,
				"invalid2",
			},
			false,
		},
		{
			"same dir",
			args{
				dir1,
				dir1,
			},
			true,
		},
		{
			"different dir (1)",
			args{
				dir1,
				dir2,
			},
			true,
		},
		{
			"different dir (2)",
			args{
				dir2,
				dir1,
			},
			true,
		},
		{
			"child dir",
			args{
				dir1Subdir,
				dir1,
			},
			false,
		},
	}
	for _, tt := range tests {
		got := isNotChildDirOf(tt.args.childPath, tt.args.parentPath)
		assert.Equal(tt.want, got, tt.name)
	}
}

func Test_isNotEmptyDir(t *testing.T) {
	assert := assert.New(t)

	dir1, err := ioutil.TempDir("", "dir")
	assert.Nil(err)
	defer os.RemoveAll(dir1)

	dir2, err := ioutil.TempDir("", "dir")
	assert.Nil(err)
	dir2File1, err := os.Create(filepath.Join(dir2, "10.txt"))
	assert.Nil(err)
	dir2File1.Close()
	defer os.RemoveAll(dir2)

	tests := []struct {
		name           string
		path           string
		includeSubdirs bool
		want           bool
	}{
		{
			"invalid dir",
			"",
			false,
			false,
		},
		{
			"empty dir (1)",
			dir1,
			false,
			false,
		},
		{
			"empty dir (2)",
			dir1,
			true,
			false,
		},
		{
			"non-empty dir (1)",
			dir2,
			false,
			true,
		},
		{
			"non-empty dir (2)",
			dir2,
			true,
			true,
		},
	}
	for _, tt := range tests {
		got := isNotEmptyDir(tt.path, tt.includeSubdirs)
		assert.Equal(tt.want, got, tt.name)
	}
}
