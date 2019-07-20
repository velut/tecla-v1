package core

// Operation represents an organizer's operation.
type Operation struct {
	ID       int64  `json:"id"`
	Op       OpType `json:"op"`
	SrcPath  string `json:"srcPath"`
	DstPath  string `json:"dstPath"`
	MaxTries int    `json:"maxTries"`
}

// OpType enum type.
type OpType string

// OpType enum values.
const (
	OpTypeCopy OpType = "copy"
	OpTypeMove OpType = "move"
)

// IsValid returns true if the OpType value belongs to the enum.
func (t OpType) IsValid() bool {
	valid := map[OpType]bool{
		OpTypeCopy: true,
		OpTypeMove: true,
	}
	return valid[t]
}
