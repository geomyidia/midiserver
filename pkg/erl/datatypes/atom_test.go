package datatypes

import (
	"testing"

	"github.com/okeuday/erlang_go/v2/erlang"
	"github.com/stretchr/testify/suite"
)

type AtomTestSuite struct {
	suite.Suite
	atom *Atom
}

func (s *AtomTestSuite) SetupTest() {
	s.atom = NewAtom("my-atom")
}

func (s *AtomTestSuite) TestValue() {
	s.Equal("my-atom", s.atom.Value())
}

func (s *AtomTestSuite) TestToTerm() {
	term, err := s.atom.ToTerm()
	s.NoError(err)
	s.Equal(erlang.OtpErlangAtom("my-atom"), term)
}

func (s *AtomTestSuite) TestFromTerm() {
	term, err := NewAtom("another-atom").ToTerm()
	s.NoError(err)
	a, err := FromTerm(term)
	s.NoError(err)
	s.Equal("another-atom", a.(*Atom).Value())
}

func TestAtomTestSuite(t *testing.T) {
	suite.Run(t, new(AtomTestSuite))
}
