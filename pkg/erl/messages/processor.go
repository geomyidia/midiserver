package messages

import (
	log "github.com/sirupsen/logrus"

	"github.com/ut-proj/midiserver/pkg/erl"
	"github.com/ut-proj/midiserver/pkg/erl/datatypes"
	"github.com/ut-proj/midiserver/pkg/erl/packets"
	"github.com/ut-proj/midiserver/pkg/types"
)

type MessageProcessor struct {
	packet    *packets.Packet
	term      interface{}
	cmdMsg    *CommandMessage
	midiCalls *MidiCallGroup
	IsMidi    bool
	IsCommand bool
}

func NewMessageProcessor(opts *erl.Opts) (*MessageProcessor, error) {
	packet, err := packets.ReadStdIOPacket(opts)
	if err != nil {
		return nil, err
	}
	t, err := packet.Term()
	if err != nil {
		return nil, err
	}
	mp := &MessageProcessor{
		packet:    packet,
		term:      t,
		IsMidi:    false,
		IsCommand: false,
	}
	log.Tracef("got Erlang Port term")
	log.Tracef("%#v", t)
	if datatypes.TupleHasKey(t, "midi") {
		mp.IsMidi = true
		callGroup, err := NewMidiCallGroup(t)
		if err != nil {
			resp := NewResponse(types.Result(""), types.Err(err.Error()))
			resp.Send()
			return nil, err
		}
		mp.midiCalls = callGroup
		return mp, nil
	}
	msg, err := NewCommandMessage(t)
	if err != nil {
		resp := NewResponse(types.Result(""), types.Err(err.Error()))
		resp.Send()
		return nil, err
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

func (mp *MessageProcessor) MidiCallGroup() *MidiCallGroup {
	return mp.midiCalls
}
