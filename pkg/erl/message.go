package erl

import (
	"errors"
	"fmt"

	erlang "github.com/okeuday/erlang_go/v2/erlang"
	log "github.com/sirupsen/logrus"

	"github.com/geomyidia/midiserver/pkg/types"
)

type CommandMessage struct {
	command erlang.OtpErlangAtom
	args    []interface{}
}

type MessageProcessor struct {
	packet  *Packet
	term    interface{}
	cmdMsg  *CommandMessage
	midiMsg interface{}
}

func handleTuple(tuple erlang.OtpErlangTuple) (*CommandMessage, error) {
	log.Debug("handling command tuple ...")
	if len(tuple) != DRCTVARITY {
		return nil, fmt.Errorf("tuple of wrong size; expected 2, got %d", len(tuple))
	}
	_, ok := tuple[DRCTVKEYINDEX].(erlang.OtpErlangAtom)
	if !ok {
		return nil, errors.New("unexpected type for directive")
	}
	msg := &CommandMessage{}
	msg.command = tuple[DRCTVVALUEINDEX].(erlang.OtpErlangAtom)
	return msg, nil
}

func handleTuples(tuples erlang.OtpErlangList) (*CommandMessage, error) {
	msg := &CommandMessage{}
	//msg.command = tuple[DRCTVVALUEINDEX].(erlang.OtpErlangAtom)
	//msg.args =
	return msg, nil
}

func NewCommandMessage(t interface{}) (*CommandMessage, error) {
	tuple, ok := t.(erlang.OtpErlangTuple)
	if !ok {
		tuples, ok := t.(erlang.OtpErlangList)
		if !ok {
			return nil, errors.New("unexpected message format")
		}
		return handleTuples(tuples)
	}
	return handleTuple(tuple)
}

func (m *CommandMessage) Command() erlang.OtpErlangAtom {
	return m.command
}

func (m *CommandMessage) Args() []interface{} {
	return m.args
}

func NewMessageProcessor(opts *Opts) (*MessageProcessor, error) {
	packet, err := ReadStdIOPacket(opts)
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
