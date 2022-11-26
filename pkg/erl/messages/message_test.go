package messages

import (
	"sort"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/suite"

	"github.com/ut-proj/midiserver/pkg/erl/datatypes"
	"github.com/ut-proj/midiserver/pkg/erl/testdata"
)

type MessageTestSuite struct {
	suite.Suite
}

func (s *MessageTestSuite) SetupSuite() {
}

func (s *MessageTestSuite) TestNewFromBytesBatch() {
	msg, err := NewFromBytes(testdata.BatchETFBytes)
	s.NoError(err)
	s.Equal("midi", msg.Type())
	s.Equal("batch", msg.Name())
	args := msg.Args()
	s.Equal(1, len(args))
	batches := args[0].(*datatypes.List).Elements()
	batch1 := batches[0].(*datatypes.Tuple).Elements()
	s.Equal("id", batch1[0].(*datatypes.Atom).Value())
	id, err := uuid.FromBytes(batch1[1].(*datatypes.Binary).Value())
	s.NoError(err)
	s.Equal("30969579-ca53-4ba0-b4af-acfced709864", id.String())
	batch2 := batches[1].(*datatypes.Tuple).Elements()
	s.Equal("messages", batch2[0].(*datatypes.Atom).Value())
	msgs := batch2[1].(*datatypes.List).Elements()
	s.Equal(4, len(msgs))
	var msgNames []string
	for _, msg := range msgs {
		msgNames = append(msgNames, msg.(*datatypes.Tuple).Key().(*datatypes.Atom).Value())
	}
	sort.Strings(msgNames)
	s.Equal([]string{"channel", "device", "note_off", "note_on"}, msgNames)
}

func TestMessageTestSuite(t *testing.T) {
	suite.Run(t, new(MessageTestSuite))
}
