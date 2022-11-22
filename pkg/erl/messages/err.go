package messages

import (
	"github.com/okeuday/erlang_go/v2/erlang"
	"github.com/ut-proj/midiserver/pkg/erl/datatypes"
	"github.com/ut-proj/midiserver/pkg/types"
)

const ErrKey = "error"

type Error struct {
	tuple *datatypes.Tuple
}

func NewError(e types.Err) *Error {
	return &Error{
		tuple: datatypes.NewTuple(ErrKey, e),
	}
}

func (e Error) Value() interface{} {
	return e.tuple.Value()
}

func (e Error) Convert() erlang.OtpErlangTuple {
	return e.tuple.Convert()
}
