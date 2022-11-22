package datatypes

import (
	"testing"

	erlang "github.com/okeuday/erlang_go/v2/erlang"
	"github.com/stretchr/testify/suite"
)

type UtilTestSuite struct {
	suite.Suite
	list erlang.OtpErlangList
}

func (s *UtilTestSuite) SetupTest() {
	list, err := NewTupleListFromSlice([][]interface{}{
		{"a", "1"},
		{"b", "2"},
	})
	s.NoError(err)
	s.list = list.Convert()
}

func (s *UtilTestSuite) TestTupleToSlice() {
	t := NewTuple("a", "1").Convert()
	slice, err := TupleToSlice(t)
	s.NoError(err)
	s.Equal(2, len(slice))
	s.Equal("a", slice[tupleKey])
	s.Equal("1", slice[tupleVal])
}

func (s *UtilTestSuite) TestTupleListToMap() {
	m, err := TupleListToMap(s.list)
	s.NoError(err)
	s.Equal("1", m["a"])
	s.Equal("2", m["b"])
}

func TestUtilTestSuite(t *testing.T) {
	suite.Run(t, new(UtilTestSuite))
}
