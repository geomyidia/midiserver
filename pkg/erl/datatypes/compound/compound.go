package compound

import (
	"github.com/okeuday/erlang_go/v2/erlang"
	errors "github.com/ut-proj/midiserver/pkg/erl/datatypes/errors"

	"github.com/ut-proj/midiserver/pkg/erl/datatypes/atom"
)

func FromTerm(term interface{}) (interface{}, error) {
	switch t := term.(type) {
	default:
		return nil, errors.ErrUnsupportedOTPType
	case erlang.OtpErlangAtom:
		return atom.FromTerm(t)
	case erlang.OtpErlangAtomCacheRef:
		return nil, errors.ErrNotImplemented
	case erlang.OtpErlangAtomUTF8:
		return nil, errors.ErrNotImplemented
	case erlang.OtpErlangBinary:
		return nil, errors.ErrNotImplemented
	case erlang.OtpErlangFunction:
		return nil, errors.ErrNotImplemented
	case erlang.OtpErlangList:
		terms := make([]interface{}, len(t.Value))
		for i, term := range t.Value {
			termStruct, err := FromTerm(term)
			if err != nil {
				return nil, errors.ErrCastingTuple
			}
			terms[i] = termStruct
		}
		return NewList(terms), nil
	case erlang.OtpErlangMap:
		return nil, errors.ErrNotImplemented
	case erlang.OtpErlangPid:
		return nil, errors.ErrNotImplemented
	case erlang.OtpErlangPort:
		return nil, errors.ErrNotImplemented
	case erlang.OtpErlangReference:
		return nil, errors.ErrNotImplemented
	case erlang.OtpErlangTuple:
		terms := make([]interface{}, len(t))
		for i, term := range t {
			termStruct, err := FromTerm(term)
			if err != nil {
				return nil, errors.ErrCastingTuple
			}
			terms[i] = termStruct
		}
		return NewTuple(terms), nil
	}
}
