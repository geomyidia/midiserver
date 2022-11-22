package datatypes

import "github.com/okeuday/erlang_go/v2/erlang"

func TupleToSlice(tuple erlang.OtpErlangTuple) ([]string, error) {
	keyTerm := tuple[tupleKey]
	key, ok := keyTerm.(erlang.OtpErlangAtom)
	if !ok {
		return nil, ErrCastingAtom
	}
	valTerm := tuple[tupleVal]
	val, ok := valTerm.(erlang.OtpErlangAtom)
	if !ok {
		return nil, ErrCastingAtom
	}
	return []string{string(key), string(val)}, nil
}

func TupleListToMap(list erlang.OtpErlangList) (map[interface{}]interface{}, error) {
	m := make(map[interface{}]interface{})
	for _, t := range list.Value {
		tpl, ok := t.(erlang.OtpErlangTuple)
		if !ok {
			return nil, ErrCastingTuple
		}
		slice, err := TupleToSlice(tpl)
		if err != nil {
			return nil, err
		}
		m[slice[tupleKey]] = slice[tupleVal]
	}
	return m, nil
}
