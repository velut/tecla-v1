package core

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestOpType_IsValid(t *testing.T) {
	assert := assert.New(t)

	tests := []struct {
		name string
		t    OpType
		want bool
	}{
		{
			"invalid op type",
			"",
			false,
		},
		{
			"valid op type copy (1)",
			"copy",
			true,
		},
		{
			"valid op type copy (2)",
			OpTypeCopy,
			true,
		},
		{
			"valid op type move (1)",
			"move",
			true,
		},
		{
			"valid op type move (2)",
			OpTypeMove,
			true,
		},
	}
	for _, tt := range tests {
		got := OpType.IsValid(tt.t)
		assert.Equal(tt.want, got, tt.name)
	}
}
