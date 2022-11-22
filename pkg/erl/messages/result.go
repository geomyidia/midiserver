package messages

import (
	"github.com/okeuday/erlang_go/v2/erlang"
	"github.com/ut-proj/midiserver/pkg/erl/datatypes"
	"github.com/ut-proj/midiserver/pkg/types"
)

const ResultKey = "result"

type Result struct {
	tuple *datatypes.Tuple
}

func NewResult(result types.Result) *Result {
	return &Result{
		tuple: datatypes.NewTuple(ResultKey, result),
	}
}

func (r Result) Value() interface{} {
	return r.tuple.Value()
}

func (r Result) Convert() erlang.OtpErlangTuple {
	return r.tuple.Convert()
}
