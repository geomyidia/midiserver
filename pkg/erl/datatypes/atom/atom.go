package atom

import (
	"github.com/okeuday/erlang_go/v2/erlang"

	"github.com/ut-proj/midiserver/pkg/erl/datatypes/errors"
)

type Atom struct {
	value string
}

func New(value string) *Atom {
	return &Atom{
		value: value,
	}
}

func (a *Atom) Value() string {
	return a.value
}

func (a *Atom) ToTerm() (interface{}, error) {
	return erlang.OtpErlangAtom(a.value), nil
}

func FromTerm(term interface{}) (*Atom, error) {
	eAtom, ok := term.(erlang.OtpErlangAtom)
	if !ok {
		return nil, errors.ErrCastingAtom
	}
	return New(string(eAtom)), nil
}
