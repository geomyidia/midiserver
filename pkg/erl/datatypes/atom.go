package datatypes

import (
	"github.com/okeuday/erlang_go/v2/erlang"
)

type Atom struct {
	value string
}

func NewAtom(value string) *Atom {
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
