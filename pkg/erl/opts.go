package erl

import (
	"github.com/geomyidia/midiserver/pkg/types"
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
