package compound

import (
	"testing"

	"github.com/okeuday/erlang_go/v2/erlang"
	"github.com/stretchr/testify/suite"

	"github.com/ut-proj/midiserver/pkg/erl/datatypes/atom"
)

type TupleTestSuite struct {
	suite.Suite
	twoTuple *Tuple
	nTuple   *Tuple
}

func (s *TupleTestSuite) SetupTest() {
	s.twoTuple = NewTuple([]interface{}{
		atom.New("a"),
		atom.New("1"),
	})
	s.nTuple = NewTuple([]interface{}{
		atom.New("b"),
		atom.New("2"),
		s.twoTuple,
		NewList([]interface{}{
			s.nTuple,
			atom.New("3"),
			atom.New("c"),
		}),
	})
}

func (s *TupleTestSuite) TestKey() {
	s.Equal("a", s.twoTuple.Key().(*atom.Atom).Value())
}

func (s *TupleTestSuite) TestValue() {
	s.Equal("1", s.twoTuple.Value().(*atom.Atom).Value())
}

func (s *TupleTestSuite) TestLen() {
	s.Equal(2, s.twoTuple.Len())
	s.Equal(4, s.nTuple.Len())
}

func (s *TupleTestSuite) TestNth() {
	s.Equal(atom.New("a"), s.twoTuple.Nth(0))
	s.Equal(atom.New("2"), s.nTuple.Nth(1))
	s.Equal(atom.New("1"), s.nTuple.Nth(2).(*Tuple).Nth(1))
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
