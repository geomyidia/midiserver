package datatypes

import (
	"testing"

	erlang "github.com/okeuday/erlang_go/v2/erlang"
	"github.com/stretchr/testify/suite"
)

type TupleTestSuite struct {
	suite.Suite
}

func (s *TupleTestSuite) SetupTest() {
}

func (s *TupleTestSuite) TestNewTupleFromTerm() {
	t1 := NewTuple("a", "1")
	t2, err := NewTupleFromTerm(t1.Convert())
	s.NoError(err)
	s.Equal(2, len(t2.elements))
	s.Equal("a", t2.Key())
	s.Equal("1", t2.Value())
}

func (s *TupleTestSuite) TestKey() {
	t := NewTuple("a", "1")
	s.Equal("a", t.Key())
}

func (s *TupleTestSuite) TestHasKey() {
	t := NewTuple("a", "1")
	s.True(t.HasKey("a"))
	s.False(t.HasKey("z"))
}

func (s *TupleTestSuite) TestValue() {
	t := NewTuple("a", "1")
	s.Equal("1", t.Value().(string))
}

func (s *TupleTestSuite) TestConvert() {
	t := NewTuple("a", "1")
	et := t.Convert()
	s.Equal(erlang.OtpErlangAtom("a"), et[tupleKey])
	s.Equal(erlang.OtpErlangAtom("1"), et[tupleVal])
}

func TestTupleTestSuite(t *testing.T) {
	suite.Run(t, new(TupleTestSuite))
}
