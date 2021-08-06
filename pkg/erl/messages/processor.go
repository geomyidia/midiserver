package messages

import (
	log "github.com/sirupsen/logrus"

	"github.com/geomyidia/midiserver/pkg/erl"
	"github.com/geomyidia/midiserver/pkg/erl/datatypes"
	"github.com/geomyidia/midiserver/pkg/erl/packets"
	"github.com/geomyidia/midiserver/pkg/types"
)

type MessageProcessor struct {
	packet    *packets.Packet
	term      interface{}
	cmdMsg    *CommandMessage
	midiCalls *MidiCalls
	IsMidi    bool
	IsCommand bool
}

func NewMessageProcessor(opts *erl.Opts) (*MessageProcessor, error) {
	nilMp := &MessageProcessor{}
	packet, err := packets.ReadStdIOPacket(opts)
	if err != nil {
		return nilMp, err
	}
	t, err := packet.Term()
	if err != nil {
		return nilMp, err
	}
	mp := &MessageProcessor{
		packet:    packet,
		term:      t,
		IsMidi:    false,
		IsCommand: false,
	}
	log.Debugf("got Erlang Port term")
	log.Tracef("%#v", t)
	if datatypes.TupleHasKey(t, "midi") {
		mp.IsMidi = true
		calls, err := NewMidiCalls(t)
		if err != nil {
			resp := NewResponse(types.Result(""), types.Err(err.Error()))
			resp.Send()
			return nilMp, err
		}
		mp.midiCalls = calls
		return mp, nil
	}
	msg, err := NewCommandMessage(t)
	if err != nil {
		resp := NewResponse(types.Result(""), types.Err(err.Error()))
		resp.Send()
		return nilMp, err
	}
	mp.cmdMsg = msg
	return mp, nil
}

func (mp *MessageProcessor) Continue() types.Result {
	return types.Result("continue")
}

func (mp *MessageProcessor) Process() types.Result {
	if mp.cmdMsg != nil {
		return types.Result(mp.cmdMsg.Command())
	} else if mp.midiCalls.Length() != 0 {
		return types.Result(types.MidiOp(types.MidiKey))
	} else {
		log.Error("unexected message type")
		return mp.Continue()
	}
}

func (mp *MessageProcessor) CommandArgs() types.PropList {
	return mp.cmdMsg.Args()
}

func (mp *MessageProcessor) MidiCalls() []types.MidiCall {
	return mp.midiCalls.Calls()
}
