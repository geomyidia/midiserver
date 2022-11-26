package datatypes

import (
	"github.com/okeuday/erlang_go/v2/erlang"
)

const (
	tupleKey   = 0
	tupleVal   = 1
	tupleArity = 2
)

type Tuple struct {
	elements []interface{}
}

func NewTuple(elements []interface{}) *Tuple {
	return &Tuple{
		elements: elements,
	}
}

func (t *Tuple) Key() interface{} {
	return t.elements[tupleKey]
}

func (t *Tuple) Value() interface{} {
	return t.elements[tupleVal]
}

func (t *Tuple) Elements() []interface{} {
	return t.elements
}

func (t *Tuple) Len() int {
	return len(t.elements)
}

func (t *Tuple) Nth(index int) interface{} {
	return t.elements[index]
}

func (t *Tuple) ToTerm() (interface{}, error) {
	terms := make([]interface{}, t.Len())
	for i, e := range t.elements {
		switch t := e.(type) {
		default:
			return nil, ErrUnsupportedGoType
		case *Atom:
			term, err := t.ToTerm()
			if err != nil {
				return nil, err
			}
			terms[i] = term
		case *Binary:
			term, err := t.ToTerm()
			if err != nil {
				return nil, err
			}
			terms[i] = term
		case *List:
			if t == nil {
				continue
			}
			term, err := t.ToTerm()
			if err != nil {
				return nil, err
			}
			terms[i] = term
		case *Tuple:
			if t == nil {
				continue
			}
			term, err := t.ToTerm()
			if err != nil {
				return nil, err
			}
			terms[i] = term
		}
	}
	return erlang.OtpErlangTuple(terms), nil
}
