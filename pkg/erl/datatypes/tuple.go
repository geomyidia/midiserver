package datatypes

import erlang "github.com/okeuday/erlang_go/v2/erlang"

const (
	tupleKey   = 0
	tupleVal   = 1
	tupleArity = 2
)

type Tuple struct {
	data []interface{}
}

func newTuple(key, value interface{}) Tuple {
	return Tuple{
		[]interface{}{key, value},
	}
}

func (t *Tuple) Key() interface{} {
	return t.data[tupleKey]
}

func (t *Tuple) Value() interface{} {
	return t.data[tupleVal]
}

func (t *Tuple) Convert() erlang.OtpErlangTuple {
	tpl := make([]interface{}, tupleArity)
	tpl[tupleKey] = erlang.OtpErlangAtom(t.Key().(string))
	tpl[tupleVal] = erlang.OtpErlangAtom(t.Value().(string))
	return erlang.OtpErlangTuple(tpl)
}
