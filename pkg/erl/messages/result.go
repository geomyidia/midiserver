package messages

import (
	"github.com/ut-proj/midiserver/pkg/erl/datatypes"
	"github.com/ut-proj/midiserver/pkg/types"
)

const ResultKey = "result"

type Result struct {
	tuple *datatypes.Tuple
}

func NewResult(result types.Result) *Result {
	return &Result{
		tuple: datatypes.NewTuple([]interface{}{
			datatypes.NewAtom(ResultKey),
			datatypes.NewAtom(string(result)),
		}),
	}
}

func (r Result) Value() interface{} {
	return r.tuple.Value()
}

func (r Result) ToTerm() (interface{}, error) {
	return r.tuple.ToTerm()
}
