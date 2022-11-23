package messages

import (
	"github.com/ut-proj/midiserver/pkg/erl/datatypes"
	"github.com/ut-proj/midiserver/pkg/types"
)

const ErrKey = "error"

type Error struct {
	tuple *datatypes.Tuple
}

func NewError(e types.Err) *Error {
	return &Error{
		tuple: datatypes.NewTuple([]interface{}{
			datatypes.NewAtom(ErrKey),
			datatypes.NewAtom(string(e)),
		}),
	}
}

func (e Error) Value() interface{} {
	return e.tuple.Value()
}

func (r Error) ToTerm() (interface{}, error) {
	return r.tuple.ToTerm()
}
