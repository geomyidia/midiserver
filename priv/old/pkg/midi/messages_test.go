package midi

import (
	"sort"
	"testing"

	"github.com/ergo-services/ergo/etf"
	"github.com/google/uuid"
	"github.com/stretchr/testify/suite"

	"github.com/geomyidia/erlcmd/pkg/messages"
	"github.com/geomyidia/erlcmd/pkg/options"
	"github.com/geomyidia/erlcmd/pkg/packets"
	"github.com/geomyidia/erlcmd/pkg/testdata"
)

const (
	Bb     uint8 = 34
	Volume uint8 = 40
)

type MidiMessageTestSuite struct {
	suite.Suite
	opts *options.Opts
}

func (s *MidiMessageTestSuite) SetupSuite() {
	s.opts = &options.Opts{IsHexEncoded: true}
}

func (s *MidiMessageTestSuite) TestBatchMessage() {
}

func (s *MidiMessageTestSuite) TestDeviceMessage() {
}

func (s *MidiMessageTestSuite) TestNotesOnMessage() {
	pkt, err := packets.NewPacket(testdata.NoteOnPacketBytes, s.opts)
	s.Require().NoError(err)
	bytes, err := pkt.Bytes()
	s.Require().NoError(err)
	msg, err := messages.NewFromBytes(bytes)
	s.Require().NoError(err)
	s.Equal("midi", msg.Type())
	s.Equal("batch", msg.Name())
	args := msg.Args()
	s.Equal(2, len(args))
	id := args[0].(etf.Tuple)
	s.Equal(etf.Atom("id"), id.Element(1).(etf.Atom))
	uuidBytes := id.Element(2).([]uint8)
	uuid, err := uuid.FromBytes(uuidBytes)
	s.Require().NoError(err)
	s.Equal("de950779-e60a-439a-bc83-327adf70d961", uuid.String())
	batch := args[1].(etf.Tuple)
	name := batch[0].(etf.Atom)
	s.Equal(etf.Atom("messages"), name)
	msgs := batch[1].(etf.List)
	s.Equal(1, len(msgs))
	var msgNames []string
	for _, msg := range msgs {
		msgNames = append(msgNames, string(msg.(etf.Tuple).Element(1).(etf.Atom)))
	}
	sort.Strings(msgNames)
	s.Equal([]string{"note_on"}, msgNames)
}

func TestMidiMessageTestSuite(t *testing.T) {
	suite.Run(t, new(MidiMessageTestSuite))
}

type MidiMessagesTestSuite struct {
	suite.Suite
	opts   *options.Opts
	batch  interface{}
	device interface{}
	noteOn interface{}
}

func (s *MidiMessagesTestSuite) SetupSuite() {
	s.opts = &options.Opts{IsHexEncoded: true}
	// Batch setup§
	bPkt, err := packets.NewPacket(testdata.BatchPacketBytes, s.opts)
	s.Require().NoError(err)
	batch, err := bPkt.ToTerm()
	s.Require().NoError(err)
	s.batch, err = messages.New(batch)
	s.Require().NoError(err)

	// Device setup
	dPkt, err := packets.NewPacket(testdata.DevicePacketBytes, s.opts)
	s.Require().NoError(err)
	device, err := dPkt.ToTerm()
	s.Require().NoError(err)
	s.device, err = messages.New(device)
	s.Require().NoError(err)

	// Note-on Setup
	nPkt, err := packets.NewPacket(testdata.NoteOnPacketBytes, s.opts)
	s.Require().NoError(err)
	noteOn, err := nPkt.ToTerm()
	s.Require().NoError(err)
	s.noteOn, err = messages.New(noteOn)
	s.Require().NoError(err)
}

// func (suite *MidiMessagesTestSuite) TestConvertDevice() {
// 	converted, err := HandleMessage(suite.device)
// 	calls := converted.(*datatypes.List).Elements()
// 	suite.Require().Equal(1, len(calls))
// 	suite.NoError(err)
// 	suite.Equal("11ff135c-78d5-415c-8818-cde72252ff02", converted.Id())
// 	suite.Equal(types.MidiDeviceType(), calls[0].Op)
// 	suite.Equal(uint8(0), calls[0].Args.Device)
// }

// func (suite *MidiMessagesTestSuite) TestConvertNoteOn() {
// 	converted, err := HandleMessage(suite.noteOn)
// 	suite.NoError(err)
// 	calls := converted.(*datatypes.List).Elements()
// 	suite.Require().Equal(1, len(calls))
// 	suite.Equal("de950779-e60a-439a-bc83-327adf70d961", converted.Id())
// 	suite.Equal(types.MidiNoteOnType(), calls[0].Op)
// 	suite.Equal(Bb, calls[0].Args.NoteOn.Pitch)
// 	suite.Equal(Volume, calls[0].Args.NoteOn.Velocity)
// }

// func (suite *MidiMessagesTestSuite) TestConvertBatch() {
// 	converted, err := HandleMessage(suite.batch)
// 	calls := converted.(*datatypes.List).Elements()
// 	suite.Require().Equal(4, len(calls))
// 	suite.NoError(err)
// 	suite.Equal("30969579-ca53-4ba0-b4af-acfced709864", converted.Id())
// 	suite.Require().Equal(4, len(calls))
// 	suite.Equal(1, calls[0].Id)
// 	suite.Equal(types.MidiDeviceType(), calls[0].Op)
// 	suite.Equal(uint8(0), calls[0].Args.Device)
// 	suite.Equal(2, calls[1].Id)
// 	suite.Equal(types.MidiChannelType(), calls[1].Op)
// 	suite.Equal(uint8(0), calls[1].Args.Channel)
// 	suite.Equal(3, calls[2].Id)
// 	suite.Equal(types.MidiNoteOnType(), calls[2].Op)
// 	suite.Equal(Bb, calls[2].Args.NoteOn.Pitch)
// 	suite.Equal(Volume, calls[2].Args.NoteOn.Velocity)
// 	suite.Equal(4, calls[3].Id)
// 	suite.Equal(types.MidiNoteOffType(), calls[3].Op)
// 	suite.Equal(Bb, calls[3].Args.NoteOff)
// }

// In order for 'go test' to run this suite, we need to create
// a normal test function and pass our suite to suite.Run
func TestMidiMessagesTestSuite(t *testing.T) {
	suite.Run(t, new(MidiMessagesTestSuite))
}
