package erl

import (
	"github.com/geomyidia/midiserver/pkg/types"
)

// Constants
const (
	DELIMITER       = '\n'
	DRCTVARITY      = 2
	DRCTVKEYINDEX   = 0
	DRCTVVALUEINDEX = 1
)

func Continue() types.Result {
	return types.Result("continue")
}

type Opts struct {
	IsHexEncoded bool
}

func DefaultOpts() *Opts {
	return &Opts{
		IsHexEncoded: false,
	}
}
