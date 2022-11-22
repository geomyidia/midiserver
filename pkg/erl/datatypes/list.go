package datatypes

import erlang "github.com/okeuday/erlang_go/v2/erlang"

type TupleList struct {
	elements []*Tuple
}

func NewTupleList(elements []*Tuple) *TupleList {
	return &TupleList{
		elements: elements,
	}
}

func NewTupleListFromSlice(pairs [][]interface{}) (*TupleList, error) {
	var elements []*Tuple
	for _, pair := range pairs {
		if len(pair) != tupleArity {
			return nil, ErrBadTupleArity
		}
		t := NewTuple(pair[0], pair[1])
		elements = append(elements, t)
	}
	return &TupleList{
		elements: elements,
	}, nil
}

func NewTupleListFromTerm(term interface{}) (*TupleList, error) {
	el, ok := term.(erlang.OtpErlangList)
	if !ok {
		return nil, ErrCastingList
	}
	var l []*Tuple
	for _, t := range el.Value {
		et, ok := t.(erlang.OtpErlangTuple)
		if !ok {
			return nil, ErrCastingTuple
		}
		tpl, err := NewTupleFromTerm(et)
		if err != nil {
			return nil, err
		}
		l = append(l, tpl)
	}
	return NewTupleList(l), nil
}

func (l *TupleList) Nth(index int) *Tuple {
	return l.elements[index]
}

func (l *TupleList) ToMap() map[interface{}]interface{} {
	m := make(map[interface{}]interface{})
	for _, t := range l.elements {
		m[t.Key()] = t.Value()
	}
	return m
}

func (l *TupleList) Convert() erlang.OtpErlangList {
	el := make([]interface{}, len(l.elements))
	for i, t := range l.elements {
		el[i] = t.Convert()
	}
	return erlang.OtpErlangList{Value: el}
}
