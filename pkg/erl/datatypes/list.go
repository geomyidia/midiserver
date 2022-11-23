package datatypes

import (
	"github.com/okeuday/erlang_go/v2/erlang"
)

type List struct {
	elements []interface{}
}

func NewList(elements []interface{}) *List {
	return &List{
		elements: elements,
	}
}

func (l *List) Len() int {
	return len(l.elements)
}

func (l *List) Nth(index int) interface{} {
	return l.elements[index]
}

func (l *List) ToTerm() (interface{}, error) {
	terms := make([]interface{}, l.Len())
	for i, e := range l.elements {
		switch t := e.(type) {
		default:
			return nil, ErrUnsupportedGoType
		case *Atom:
			term, err := t.ToTerm()
			if err != nil {
				return nil, err
			}
			terms[i] = term
		case *List:
			term, err := t.ToTerm()
			if err != nil {
				return nil, err
			}
			terms[i] = term
		case *Tuple:
			term, err := t.ToTerm()
			if err != nil {
				return nil, err
			}
			terms[i] = term
		}
	}
	return erlang.OtpErlangList{Value: terms}, nil
}

// type TupleList struct {
// 	elements []*Tuple
// }

// func NewTupleList(elements []*Tuple) *TupleList {
// 	return &TupleList{
// 		elements: elements,
// 	}
// }

// func NewTupleListFromSlice(pairs [][]interface{}) (*TupleList, error) {
// 	var elements []*Tuple
// 	for _, pair := range pairs {
// 		if len(pair) != tupleArity {
// 			return nil, ErrBadTupleArity
// 		}
// 		t := NewTuple(pair[0].(string), pair[1])
// 		elements = append(elements, t)
// 	}
// 	return &TupleList{
// 		elements: elements,
// 	}, nil
// }

// func NewTupleListFromTerm(term interface{}) (*TupleList, error) {
// 	el, ok := term.(erlang.OtpErlangList)
// 	if !ok {
// 		return nil, ErrCastingList
// 	}
// 	var l []*Tuple
// 	for _, t := range el.Value {
// 		et, ok := t.(erlang.OtpErlangTuple)
// 		if !ok {
// 			return nil, ErrCastingTuple
// 		}
// 		tpl, err := NewTupleFromTerm(et)
// 		if err != nil {
// 			return nil, err
// 		}
// 		l = append(l, tpl)
// 	}
// 	return NewTupleList(l), nil
// }

// func (l *TupleList) Nth(index int) *Tuple {
// 	return l.elements[index]
// }

// func (l *TupleList) ToMap() map[interface{}]interface{} {
// 	m := make(map[interface{}]interface{})
// 	for _, t := range l.elements {
// 		m[t.Key()] = t.Value()
// 	}
// 	return m
// }

// func (l *TupleList) Convert() erlang.OtpErlangList {
// 	el := make([]interface{}, len(l.elements))
// 	for i, t := range l.elements {
// 		el[i] = t.Convert()
// 	}
// 	return erlang.OtpErlangList{Value: el}
// }
