package erl

import (
	"github.com/ut-proj/midiserver/pkg/types"
)

func Continue() types.Result {
	return types.ContinueResult
}

type Opts struct {
	IsHexEncoded bool
}

func DefaultOpts() *Opts {
	return &Opts{
		IsHexEncoded: false,
	}
}
