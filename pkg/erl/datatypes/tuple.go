package datatypes

import erlang "github.com/okeuday/erlang_go/v2/erlang"

const (
	tupleKey   = 0
	tupleVal   = 1
	tupleArity = 2
)

type Tuple struct {
	elements []interface{}
}

func NewTuple(key, value interface{}) *Tuple {
	return &Tuple{
		elements: []interface{}{key, value},
	}
}

func NewTupleFromTerm(term interface{}) (*Tuple, error) {
	et, ok := term.(erlang.OtpErlangTuple)
	if !ok {
		return nil, ErrCastingTuple
	}
	slice, err := TupleToSlice(et)
	if err != nil {
		return nil, err
	}
	return NewTuple(slice[tupleKey], slice[tupleVal]), nil
}

func (t *Tuple) Key() interface{} {
	return t.elements[tupleKey]
}

func (t *Tuple) Value() interface{} {
	return t.elements[tupleVal]
}

func (t *Tuple) HasKey(key interface{}) bool {
	return t.Key() == key
}

func (t *Tuple) Convert() erlang.OtpErlangTuple {
	tpl := make([]interface{}, tupleArity)
	tpl[tupleKey] = erlang.OtpErlangAtom(t.Key().(string))
	tpl[tupleVal] = erlang.OtpErlangAtom(t.Value().(string))
	return erlang.OtpErlangTuple(tpl)
}
