package datatypes

import (
	"testing"

	"github.com/okeuday/erlang_go/v2/erlang"
	"github.com/stretchr/testify/suite"
)

type BinaryTestSuite struct {
	suite.Suite
	binary *Binary
}

func (s *BinaryTestSuite) SetupTest() {
	s.binary = NewBinary([]byte("my binary"))
}

func (s *BinaryTestSuite) TestFromStr() {
	str := "hey, a string! but it's a binary! but it's a string!"
	b := NewBinaryFromStr(str)
	s.Equal(str, string(b.Value()))
}

func (s *BinaryTestSuite) TestValue() {
	s.Equal("my binary", string(s.binary.Value()))
}

func (s *BinaryTestSuite) TestToTerm() {
	term, err := s.binary.ToTerm()
	s.NoError(err)
	s.Equal(erlang.OtpErlangBinary{Value: []byte("my binary")}, term)
}

func (s *BinaryTestSuite) TestFromTerm() {
	term, err := NewBinaryFromStr("another binary").ToTerm()
	s.NoError(err)
	a, err := FromTerm(term)
	s.NoError(err)
	s.Equal("another binary", string(a.(*Binary).Value()))
}

func TestBinaryTestSuite(t *testing.T) {
	suite.Run(t, new(BinaryTestSuite))
}
