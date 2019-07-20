package core

import (
	"encoding/json"
	"fmt"
	"strings"
	"unicode/utf8"

	"github.com/velut/fsutils-go/fs"
)

// Error keys and error messages.
const (
	// Config keys
	ErrKeyConfig     = "config"
	ErrKeyConfigName = "config.name"
	// Config errors
	ErrConfigNil       = "no configuration found"
	ErrConfigNameEmpty = "name is empty"

	// ConfigSrc keys
	ErrKeySrc              = "config.src"
	ErrKeySrcDir           = "config.src.dir"
	ErrKeySrcDefaultOpType = "config.src.defaultOpType"
	// ConfigSrc errors
	ErrSrcNil                   = "no configuration found"
	ErrSrcDirPathEmpty          = "path is empty"
	ErrSrcDirPathNotValid       = "path is not valid"
	ErrSrcDirEmpty              = "directory contains no files"
	ErrSrcDefaultOpTypeNotValid = "default operation type is not valid"

	// ConfigDst keys
	ErrKeyDst          = "config.dst"
	ErrKeyDstDirs      = "config.dst.dirs"
	ErrKeyDstDirHotkey = "config.dst.dirs.%v.hotkey"
	ErrKeyDstDirPath   = "config.dst.dirs.%v.dir"
	// ConfigDst errors
	ErrDstNil                           = "no configuration found"
	ErrDstDirsEmpty                     = "no destination directories"
	ErrDstDirHotkeyEmpty                = "hotkey is empty"
	ErrDstDirHotkeyNotOneRune           = "hotkey is too long"
	ErrDstDirHotkeyDuplicate            = "hotkey is a duplicate"
	ErrDstDirPathEmpty                  = "path is empty"
	ErrDstDirPathNotValid               = "path is not valid"
	ErrDstDirPathNotDifferentFromSrcDir = "path points to the source directory"
	ErrDstDirPathChildOfSrcDir          = "path is inside the source directory"

	// ConfigOps keys
	ErrKeyOps           = "config.ops"
	ErrKeyOpsNumWorkers = "config.ops.numWorkers"
	ErrKeyOpsMaxTries   = "config.ops.maxTries"
	// ConfigOps errors
	ErrOpsNil                        = "no configuration found"
	ErrOpsNumWorkersNotAtLeastOne    = "number of workers is less than one"
	ErrOpsNumWorkersMoreThanFive     = "number of workers is more than five"
	ErrOpsMaxTriesNotAtLeastOne      = "number of maximum operation tries is less than one"
	ErrOpsMaxTriesMoreThanOneMillion = "number of maximum operation tries is more than one million"
)

// ConfigValidator represents the validator for configurations.
type ConfigValidator struct{}

// configValidator is the real validator used by ConfigValidator.
type configValidator struct {
	config *Config
	errs   *ConfigValidationError
}

// ConfigValidationError contains the validation errors.
type ConfigValidationError struct {
	Errors errorsByKey `json:"errors"`
}

// errorsByKey represents the mapping from error keys to error messages.
type errorsByKey map[string]string

// Error implements the error interface.
func (e *ConfigValidationError) Error() string {
	j, _ := json.Marshal(e)
	return string(j)
}

// NewConfigValidator returns a new ConfigValidator.
func NewConfigValidator() *ConfigValidator {
	return &ConfigValidator{}
}

// ValidateConfig validates the given configuration,
// returning an error of type ConfigValidationError if the configuration is not valid.
func (cv *ConfigValidator) ValidateConfig(config *Config) error {
	errs := &ConfigValidationError{make(errorsByKey)}
	validator := &configValidator{config, errs}
	return validator.validate()
}

func (v *configValidator) validate() error {
	if v.anyError() {
		return v.errs
	}
	return nil
}

func (v *configValidator) anyError() bool {
	for _, p := range v.validationPredicates() {
		if !p() {
			return true
		}
	}
	return false
}

func (v *configValidator) validationPredicates() []func() bool {
	return []func() bool{
		// Config
		v.isConfigNotNil,
		v.isSrcNotNil,
		v.isDstNotNil,
		v.isOpsNotNil,
		v.isConfigNameNotEmpty,
		// ConfigSrc
		v.isSrcDirPathNotEmpty,
		v.isSrcDirPathValid,
		v.isSrcDirNotEmpty,
		v.isSrcDefaultOpTypeValid,
		// ConfigDst
		v.areDstDirsNotEmpty,
		v.areDstDirsHotkeysAllNotEmpty,
		v.areDstDirsHotkeysAllOneRune,
		v.areDstDirsHotkeysAllDistinct,
		v.areDstDirsPathsAllNotEmpty,
		v.areDstDirsPathsAllValid,
		v.areDstDirsPathsAllDifferentFromSrcDir,
		v.areDstDirsPathsAllNotChildrenOfSrcDir,
		// ConfigOps
		v.isOpsNumWorkersAtLeastOne,
		v.isOpsNumWorkersLessThanFive,
		v.isOpsMaxTriesAtLeastOne,
		v.isOpsMaxTriesLessThanOneMillion,
	}
}

func (v *configValidator) isConfigNotNil() bool {
	ok := v.config != nil
	v.addErrIf(!ok, ErrKeyConfig, ErrConfigNil)
	return ok
}

func (v *configValidator) isSrcNotNil() bool {
	ok := v.config.Src != nil
	v.addErrIf(!ok, ErrKeySrc, ErrSrcNil)
	return ok
}

func (v *configValidator) isDstNotNil() bool {
	ok := v.config.Dst != nil
	v.addErrIf(!ok, ErrKeyDst, ErrDstNil)
	return ok
}

func (v *configValidator) isOpsNotNil() bool {
	ok := v.config.Ops != nil
	v.addErrIf(!ok, ErrKeyOps, ErrOpsNil)
	return ok
}

func (v *configValidator) isConfigNameNotEmpty() bool {
	ok := isNotEmptyString(v.config.Name)
	v.addErrIf(!ok, ErrKeyConfigName, ErrConfigNameEmpty)
	return ok
}

func (v *configValidator) isSrcDirPathNotEmpty() bool {
	ok := isNotEmptyString(v.config.Src.Dir)
	v.addErrIf(!ok, ErrKeySrcDir, ErrSrcDirPathEmpty)
	return ok
}

func (v *configValidator) isSrcDirPathValid() bool {
	ok := isDir(v.config.Src.Dir)
	v.addErrIf(!ok, ErrKeySrcDir, ErrSrcDirPathNotValid)
	return ok
}

func (v *configValidator) isSrcDirNotEmpty() bool {
	ok := isNotEmptyDir(v.config.Src.Dir, v.config.Src.IncludeSubdirs)
	v.addErrIf(!ok, ErrKeySrcDir, ErrSrcDirEmpty)
	return ok
}

func (v *configValidator) isSrcDefaultOpTypeValid() bool {
	ok := v.config.Src.DefaultOpType.IsValid()
	v.addErrIf(!ok, ErrKeySrcDefaultOpType, ErrSrcDefaultOpTypeNotValid)
	return ok
}

func (v *configValidator) areDstDirsNotEmpty() bool {
	ok := len(v.config.Dst.Dirs) > 0
	v.addErrIf(!ok, ErrKeyDstDirs, ErrDstDirsEmpty)
	return ok
}

func (v *configValidator) areDstDirsHotkeysAllNotEmpty() bool {
	allOk := true
	dirs := v.config.Dst.Dirs
	for i, d := range dirs {
		ok := isNotEmptyString(d.Hotkey)
		v.addErrWithIndexIf(!ok, ErrKeyDstDirHotkey, i, ErrDstDirHotkeyEmpty)
		allOk = allOk && ok
	}
	return allOk
}

func (v *configValidator) areDstDirsHotkeysAllOneRune() bool {
	allOk := true
	dirs := v.config.Dst.Dirs
	for i, d := range dirs {
		ok := utf8.RuneCountInString(d.Hotkey) == 1
		v.addErrWithIndexIf(!ok, ErrKeyDstDirHotkey, i, ErrDstDirHotkeyNotOneRune)
		allOk = allOk && ok
	}
	return allOk
}

func (v *configValidator) areDstDirsHotkeysAllDistinct() bool {
	allOk := true
	dirs := v.config.Dst.Dirs
	seen := make(map[string]int)
	for i, d := range dirs {
		seen[d.Hotkey]++
		ok := seen[d.Hotkey] == 1
		v.addErrWithIndexIf(!ok, ErrKeyDstDirHotkey, i, ErrDstDirHotkeyDuplicate)
		allOk = allOk && ok
	}
	return allOk
}

func (v *configValidator) areDstDirsPathsAllNotEmpty() bool {
	allOk := true
	dirs := v.config.Dst.Dirs
	for i, d := range dirs {
		ok := isNotEmptyString(d.Dir)
		v.addErrWithIndexIf(!ok, ErrKeyDstDirPath, i, ErrDstDirPathEmpty)
		allOk = allOk && ok
	}
	return allOk
}

func (v *configValidator) areDstDirsPathsAllValid() bool {
	allOk := true
	dirs := v.config.Dst.Dirs
	for i, d := range dirs {
		ok := isDir(d.Dir)
		v.addErrWithIndexIf(!ok, ErrKeyDstDirPath, i, ErrDstDirPathNotValid)
		allOk = allOk && ok
	}
	return allOk
}

func (v *configValidator) areDstDirsPathsAllDifferentFromSrcDir() bool {
	allOk := true
	srcDir := v.config.Src.Dir
	dirs := v.config.Dst.Dirs
	for i, d := range dirs {
		ok := areNotSameDir(d.Dir, srcDir)
		v.addErrWithIndexIf(!ok, ErrKeyDstDirPath, i, ErrDstDirPathNotDifferentFromSrcDir)
		allOk = allOk && ok
	}
	return allOk
}

func (v *configValidator) areDstDirsPathsAllNotChildrenOfSrcDir() bool {
	allOk := true
	srcDir := v.config.Src.Dir
	dirs := v.config.Dst.Dirs
	for i, d := range dirs {
		ok := isNotChildDirOf(d.Dir, srcDir)
		v.addErrWithIndexIf(!ok, ErrKeyDstDirPath, i, ErrDstDirPathChildOfSrcDir)
		allOk = allOk && ok
	}
	return allOk
}

func (v *configValidator) isOpsNumWorkersAtLeastOne() bool {
	ok := v.config.Ops.NumWorkers >= 1
	v.addErrIf(!ok, ErrKeyOpsNumWorkers, ErrOpsNumWorkersNotAtLeastOne)
	return ok
}

func (v *configValidator) isOpsNumWorkersLessThanFive() bool {
	ok := v.config.Ops.NumWorkers <= 5
	v.addErrIf(!ok, ErrKeyOpsNumWorkers, ErrOpsNumWorkersMoreThanFive)
	return ok
}

func (v *configValidator) isOpsMaxTriesAtLeastOne() bool {
	ok := v.config.Ops.MaxTries >= 1
	v.addErrIf(!ok, ErrKeyOpsMaxTries, ErrOpsMaxTriesNotAtLeastOne)
	return ok
}

func (v *configValidator) isOpsMaxTriesLessThanOneMillion() bool {
	ok := v.config.Ops.MaxTries <= 1000000
	v.addErrIf(!ok, ErrKeyOpsMaxTries, ErrOpsMaxTriesMoreThanOneMillion)
	return ok
}

func (v *configValidator) addErrWithIndexIf(add bool, keyFmt string, index int, val string) {
	key := fmt.Sprintf(keyFmt, index)
	v.addErrIf(add, key, val)
}

func (v *configValidator) addErrIf(add bool, key, val string) {
	if add {
		v.errs.Add(key, val)
	}
}

// Add adds an error to the validation errors.
func (e *ConfigValidationError) Add(key, val string) {
	e.Errors[key] = val
}

func isNotEmptyString(s string) bool {
	return len(strings.TrimSpace(s)) > 0
}

func isDir(path string) bool {
	err := fs.AssertDir(path)
	return err == nil
}

func isNotEmptyDir(path string, includeSubdirs bool) bool {
	fis, err := fs.ReadDir(path, &fs.ReadDirOptions{
		IncludeSubdirs: includeSubdirs,
		MaxFiles:       1,
	})
	if err != nil {
		return false
	}
	return len(fis) >= 1
}

func areNotSameDir(path1, path2 string) bool {
	isSame, err := fs.SameDir(path1, path2)
	if err != nil {
		return false
	}
	return !isSame
}

func isNotChildDirOf(childPath, parentPath string) bool {
	isChild, err := fs.SubdirOf(childPath, parentPath)
	if err != nil {
		return false
	}
	return !isChild
}
