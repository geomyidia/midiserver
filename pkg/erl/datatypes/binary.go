package datatypes

import (
	"github.com/okeuday/erlang_go/v2/erlang"
)

type Binary struct {
	value []byte
}

func NewBinary(value []byte) *Binary {
	return &Binary{
		value: value,
	}
}

func NewBinaryFromStr(value string) *Binary {
	return NewBinary([]byte(value))
}

func (b *Binary) Value() []byte {
	return b.value
}

func (b *Binary) ToTerm() (interface{}, error) {
	return erlang.OtpErlangBinary{Value: b.value}, nil
}
