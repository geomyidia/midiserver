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

func (s *TupleTestSuite) TestKey() {
	t := newTuple("a", "1")
	s.Equal("a", t.Key().(string))
}

func (s *TupleTestSuite) TestValue() {
	t := newTuple("a", "1")
	s.Equal("1", t.Value().(string))
}

func (s *TupleTestSuite) TestConvert() {
	t := newTuple("a", "1")
	et := t.Convert()
	s.Equal(erlang.OtpErlangAtom("a"), et[tupleKey])
	s.Equal(erlang.OtpErlangAtom("1"), et[tupleVal])
}

func TestTupleTestSuite(t *testing.T) {
	suite.Run(t, new(TupleTestSuite))
}
