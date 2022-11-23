package datatypes

import (
	"github.com/okeuday/erlang_go/v2/erlang"
)

func FromTerm(term interface{}) (interface{}, error) {
	switch t := term.(type) {
	default:
		return nil, ErrUnsupportedOTPType
	case erlang.OtpErlangAtom:
		eAtom, ok := term.(erlang.OtpErlangAtom)
		if !ok {
			return nil, ErrCastingAtom
		}
		return NewAtom(string(eAtom)), nil
	case erlang.OtpErlangAtomCacheRef:
		return nil, ErrNotImplemented
	case erlang.OtpErlangAtomUTF8:
		return nil, ErrNotImplemented
	case erlang.OtpErlangBinary:
		return nil, ErrNotImplemented
	case erlang.OtpErlangFunction:
		return nil, ErrNotImplemented
	case erlang.OtpErlangList:
		terms := make([]interface{}, len(t.Value))
		for i, term := range t.Value {
			termStruct, err := FromTerm(term)
			if err != nil {
				return nil, ErrCastingTuple
			}
			terms[i] = termStruct
		}
		return NewList(terms), nil
	case erlang.OtpErlangMap:
		return nil, ErrNotImplemented
	case erlang.OtpErlangPid:
		return nil, ErrNotImplemented
	case erlang.OtpErlangPort:
		return nil, ErrNotImplemented
	case erlang.OtpErlangReference:
		return nil, ErrNotImplemented
	case erlang.OtpErlangTuple:
		terms := make([]interface{}, len(t))
		for i, term := range t {
			termStruct, err := FromTerm(term)
			if err != nil {
				return nil, ErrCastingTuple
			}
			terms[i] = termStruct
		}
		return NewTuple(terms), nil
	}
}
