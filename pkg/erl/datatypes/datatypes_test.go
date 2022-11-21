package datatypes

import (
	"testing"

	erlang "github.com/okeuday/erlang_go/v2/erlang"
	"github.com/stretchr/testify/suite"
)

type DatatypesTestSuite struct {
	suite.Suite
	propList erlang.OtpErlangList
	tuple    erlang.OtpErlangTuple
}

func (s *DatatypesTestSuite) SetupTest() {
	t1 := make([]interface{}, 2)
	t1[0] = erlang.OtpErlangAtom("a")
	t1[1] = erlang.OtpErlangAtom("1")
	et1 := erlang.OtpErlangTuple(t1)
	s.tuple = et1
	t2 := make([]interface{}, 2)
	t2[0] = erlang.OtpErlangAtom("b")
	t2[1] = erlang.OtpErlangAtom("2")
	et2 := erlang.OtpErlangTuple(t2)
	s.propList = erlang.OtpErlangList{Value: []interface{}{et1, et2}}
}

func (s *DatatypesTestSuite) TestPropListToMap() {
	data, err := PropListToMap(s.propList)
	s.NoError(err)
	s.Equal(erlang.OtpErlangAtom("1"), data["a"])
	s.Equal(erlang.OtpErlangAtom("2"), data["b"])
}

func (s *DatatypesTestSuite) TestTuple() {
	key, val, err := Tuple(s.tuple)
	s.NoError(err)
	s.Equal("a", key)
	s.Equal(erlang.OtpErlangAtom("1"), val)
}

func (s *DatatypesTestSuite) TestTupleHasKey() {
	s.True(TupleHasKey(s.tuple, "a"))
	s.False(TupleHasKey(s.tuple, "b"))
}

func TestDatatypesTestSuite(t *testing.T) {
	suite.Run(t, new(DatatypesTestSuite))
}
