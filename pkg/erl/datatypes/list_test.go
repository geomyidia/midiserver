package datatypes

import (
	"testing"

	erlang "github.com/okeuday/erlang_go/v2/erlang"
	"github.com/stretchr/testify/suite"
)

type ListTestSuite struct {
	suite.Suite
	list *TupleList
}

func (s *ListTestSuite) SetupTest() {
	var err error
	s.list, err = NewTupleListFromSlice([][]interface{}{
		{"a", "1"},
		{"b", "2"},
	})
	s.NoError(err)
	s.Equal(2, len(s.list.elements))
}

func (s *ListTestSuite) TestNewTupleListFromTerm() {
	list, err := NewTupleListFromTerm(s.list.Convert())
	s.NoError(err)
	s.Equal(2, len(list.elements))
	m := list.ToMap()
	s.Equal("1", m["a"])
	s.Equal("2", m["b"])
}

func (s *ListTestSuite) TestNth() {
	t := s.list.Nth(1)
	s.Equal("b", t.Key())
	s.Equal("2", t.Value())
}

func (s *ListTestSuite) TestToMap() {
	m := s.list.ToMap()
	s.Equal("1", m["a"])
	s.Equal("2", m["b"])
}

func (s *ListTestSuite) TestConvert() {
	el := s.list.Convert()
	s.Equal(2, len(el.Value))
	et := el.Value[0].(erlang.OtpErlangTuple)
	s.Equal(erlang.OtpErlangAtom("a"), et[tupleKey])
	s.Equal(erlang.OtpErlangAtom("1"), et[tupleVal])
}

func TestListTestSuite(t *testing.T) {
	suite.Run(t, new(ListTestSuite))
}
