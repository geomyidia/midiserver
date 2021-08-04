package messages

import (
	log "github.com/sirupsen/logrus"

	"github.com/geomyidia/midiserver/pkg/erl"
	"github.com/geomyidia/midiserver/pkg/erl/packets"
	"github.com/geomyidia/midiserver/pkg/types"
)

type MessageProcessor struct {
	packet  *packets.Packet
	term    interface{}
	cmdMsg  *CommandMessage
	midiMsg interface{}
}

func NewMessageProcessor(opts *erl.Opts) (*MessageProcessor, error) {
	packet, err := packets.ReadStdIOPacket(opts)
	if err != nil {
		return &MessageProcessor{}, err
	}
	t, err := packet.Term()
	if err != nil {
		return &MessageProcessor{}, err
	}
	log.Debugf("got Erlang Port term")
	log.Tracef("%#v", t)
	msg, err := NewCommandMessage(t)
	if err != nil {
		resp := NewResponse(types.Result(""), types.Err(err.Error()))
		resp.Send()
		return &MessageProcessor{}, err
	}
	return &MessageProcessor{
		packet: packet,
		term:   t,
		cmdMsg: msg,
	}, nil
}

func (mp *MessageProcessor) Continue() types.Result {
	return types.Result("continue")
}

func (mp *MessageProcessor) Process() types.Result {
	if mp.cmdMsg != nil {
		return types.Result(mp.cmdMsg.Command())
	} else if mp.midiMsg != nil {
		// process MIDI message
		return mp.Continue()
	} else {
		log.Error("unexected message type")
		return mp.Continue()
	}
}

func (mp *MessageProcessor) CommandArgs() types.Proplist {
	return mp.cmdMsg.Args()
}

func (mp *MessageProcessor) Midi() interface{} {
	return mp.midiMsg
}
