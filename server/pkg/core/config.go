package core

// Config represents a configuration usable by the organizer.
type Config struct {
	ID   int64      `json:"id"`
	Name string     `json:"name"`
	Src  *ConfigSrc `json:"src"`
	Dst  *ConfigDst `json:"dst"`
	Ops  *ConfigOps `json:"ops"`
}

// ConfigSrc contains the configuration options for the source directory.
type ConfigSrc struct {
	Dir            string `json:"dir"`
	IncludeSubdirs bool   `json:"includeSubdirs"`
	DefaultOpType  OpType `json:"defaultOpType"`
}

// ConfigDst contains the configuration options for the destination directories.
type ConfigDst struct {
	Dirs []*DstDir `json:"dirs"`
}

// DstDir represents a destination directory.
type DstDir struct {
	Hotkey string `json:"hotkey"`
	Dir    string `json:"dir"`
}

// ConfigOps contains the configuration options specifying how operations should be executed.
type ConfigOps struct {
	NumWorkers int `json:"numWorkers"`
	MaxTries   int `json:"maxTries"`
}
