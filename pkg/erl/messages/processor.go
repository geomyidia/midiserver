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
	midiMsg   *MidiMessage
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
		msg, err := NewMidiMessage(t)
		if err != nil {
			resp := NewResponse(types.Result(""), types.Err(err.Error()))
			resp.Send()
			return nilMp, err
		}
		mp.midiMsg = msg
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
	} else if mp.midiMsg != nil {
		return types.Result(mp.midiMsg.Op())
	} else {
		log.Error("unexected message type")
		return mp.Continue()
	}
}

func (mp *MessageProcessor) CommandArgs() types.PropList {
	return mp.cmdMsg.Args()
}

func (mp *MessageProcessor) Midi() *MidiMessage {
	return mp.midiMsg
}

func (mp *MessageProcessor) MidiOp() types.MidiOpType {
	return mp.midiMsg.Op()
}

func (mp *MessageProcessor) MidiData() interface{} {
	return mp.midiMsg.Data()
}
