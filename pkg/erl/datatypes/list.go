package datatypes

import erlang "github.com/okeuday/erlang_go/v2/erlang"

type TupleList struct {
	elements []Tuple
}

func newTupleList(elements []Tuple) TupleList {
	return TupleList{
		elements: elements,
	}
}

func (l *TupleList) Nth(index int) Tuple {
	return l.elements[index]
}

func (l *TupleList) Convert() erlang.OtpErlangList {
	el := make([]interface{}, len(l.elements))
	for i, t := range l.elements {
		el[i] = t.Convert()
	}
	return erlang.OtpErlangList{Value: el}
}
