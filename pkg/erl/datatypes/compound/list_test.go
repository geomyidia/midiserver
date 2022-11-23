package compound

import (
	"testing"

	"github.com/okeuday/erlang_go/v2/erlang"
	"github.com/stretchr/testify/suite"
	"github.com/ut-proj/midiserver/pkg/erl/datatypes/atom"
)

type ListTestSuite struct {
	suite.Suite
	simpleList *List
	list       *List
}

func (s *ListTestSuite) SetupTest() {
	s.simpleList = NewList([]interface{}{
		atom.New("1"),
		atom.New("2"),
		atom.New("3"),
	})
	s.list = NewList([]interface{}{
		s.simpleList,
		atom.New("hydrogen"),
		NewTuple([]interface{}{
			atom.New("never"),
			atom.New("let"),
			atom.New("you"),
			atom.New("down"),
		}),
		atom.New("42"),
	})
}

func (s *ListTestSuite) TestLen() {
	s.Equal(3, s.simpleList.Len())
	s.Equal(4, s.list.Len())
}

func (s *ListTestSuite) TestNth() {
	s.Equal(atom.New("2"), s.simpleList.Nth(1))
	s.Equal(atom.New("never"), s.list.Nth(2).(*Tuple).Nth(0))
	s.Equal(atom.New("3"), s.list.Nth(0).(*List).Nth(2))
}

func (s *ListTestSuite) TestToTerm() {
	term, err := s.simpleList.ToTerm()
	s.NoError(err)
	l1, ok := term.(erlang.OtpErlangList)
	s.True(ok)
	s.Equal(3, len(l1.Value))
	s.Equal(erlang.OtpErlangAtom("2"), l1.Value[1])
	term, err = s.list.ToTerm()
	s.NoError(err)
	l2, ok := term.(erlang.OtpErlangList)
	s.True(ok)
	s.Equal(4, len(l2.Value))
	innerTuple, ok := l2.Value[2].(erlang.OtpErlangTuple)
	s.True(ok)
	s.Equal(erlang.OtpErlangAtom("down"), innerTuple[3])
}

func (s *ListTestSuite) TestFromTerm() {
	term, err := s.simpleList.ToTerm()
	s.NoError(err)
	data, err := FromTerm(term)
	s.NoError(err)
	simpleList, ok := data.(*List)
	s.True(ok)
	s.Equal(atom.New("3"), simpleList.Nth(2).(*atom.Atom))
	term, err = s.list.ToTerm()
	s.NoError(err)
	data, err = FromTerm(term)
	s.NoError(err)
	list, ok := data.(*List)
	s.True(ok)
	s.Equal(atom.New("never"), list.Nth(2).(*Tuple).Nth(0))
}

func TestListTestSuite(t *testing.T) {
	suite.Run(t, new(ListTestSuite))
}
