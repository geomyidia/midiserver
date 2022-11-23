package datatypes

import (
	"testing"

	"github.com/okeuday/erlang_go/v2/erlang"
	"github.com/stretchr/testify/suite"
)

type TupleTestSuite struct {
	suite.Suite
	twoTuple *Tuple
	nTuple   *Tuple
}

func (s *TupleTestSuite) SetupTest() {
	s.twoTuple = NewTuple([]interface{}{
		NewAtom("a"),
		NewAtom("1"),
	})
	s.nTuple = NewTuple([]interface{}{
		NewAtom("b"),
		NewAtom("2"),
		s.twoTuple,
		NewList([]interface{}{
			s.nTuple,
			NewAtom("3"),
			NewAtom("c"),
		}),
	})
}

func (s *TupleTestSuite) TestKey() {
	s.Equal("a", s.twoTuple.Key().(*Atom).Value())
}

func (s *TupleTestSuite) TestValue() {
	s.Equal("1", s.twoTuple.Value().(*Atom).Value())
}

func (s *TupleTestSuite) TestLen() {
	s.Equal(2, s.twoTuple.Len())
	s.Equal(4, s.nTuple.Len())
}

func (s *TupleTestSuite) TestNth() {
	s.Equal(NewAtom("a"), s.twoTuple.Nth(0))
	s.Equal(NewAtom("2"), s.nTuple.Nth(1))
	s.Equal(NewAtom("1"), s.nTuple.Nth(2).(*Tuple).Nth(1))
}

func (s *TupleTestSuite) TestToTerm() {
	term, err := s.nTuple.ToTerm()
	s.NoError(err)
	tuple, ok := term.(erlang.OtpErlangTuple)
	s.True(ok)
	s.Equal(erlang.OtpErlangAtom("b"), tuple[0])
	s.Equal(erlang.OtpErlangAtom("2"), tuple[1])
	innerTuple, ok := tuple[2].(erlang.OtpErlangTuple)
	s.True(ok)
	s.Equal(erlang.OtpErlangAtom("a"), innerTuple[0])
	s.Equal(erlang.OtpErlangAtom("1"), innerTuple[1])
}

func TestTupleTestSuite(t *testing.T) {
	suite.Run(t, new(TupleTestSuite))
}
