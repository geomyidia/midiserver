package erl

import (
	"errors"

	erlang "github.com/okeuday/erlang_go/v2/erlang"
	log "github.com/sirupsen/logrus"

	"github.com/geomyidia/midiserver/pkg/erl/proplists"
	"github.com/geomyidia/midiserver/pkg/types"
)

type CommandMessage struct {
	command types.CommandType
	args    types.Proplist
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

func (cm *CommandMessage) Command() types.CommandType {
	return cm.command
}

func (cm *CommandMessage) Args() types.Proplist {
	return cm.args
}

func (cm *CommandMessage) SetCommand(cmdIf interface{}) error {
	cmdAtom, ok := cmdIf.(erlang.OtpErlangAtom)
	if ! ok {
		return errors.New("could not cast command to atom")
	}
	cm.command = types.Command(types.CommandName(string(cmdAtom)))
	return nil
}

func (cm *CommandMessage) SetArgs(argsIf interface{}) error {
	args, err := proplists.ToMap(argsIf.(erlang.OtpErlangList))
	if err != nil {
		return err
	}
	cm.args = args
	return nil
}

type MessageProcessor struct {
	packet  *Packet
	term    interface{}
	cmdMsg  *CommandMessage
	midiMsg interface{}
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

func (mp *MessageProcessor) CommandArgs() types.Proplist {
	return mp.cmdMsg.Args()
}

func (mp *MessageProcessor) Midi() interface{} {
	return mp.midiMsg
}

func handleTuple(tuple erlang.OtpErlangTuple) (*CommandMessage, error) {
	log.Debug("handling tuple ...")
	key, val, err := proplists.ExtractTuple(tuple)
	if err != nil {
		log.Error(err)
		return nil, err
	}
	log.Debugf("Key: %+v (type %T)", key, key)
	if key == types.CommandKey {
		msg := &CommandMessage{}
		err = msg.SetCommand(val)
		if err != nil {
			log.Error(err)
			return nil, err
		}
		return msg, nil
	} else if key == types.MidiKey {
		log.Warning("MIDI Support: TODO")
		return nil, nil
	}
	return nil, nil
}

func handleTuples(tuples erlang.OtpErlangList) (*CommandMessage, error) {
	log.Debug("handling tuples ...")
	t, err := proplists.ToMap(tuples)
	if err != nil {
		log.Error(err)
		return nil, err
	}
	log.Debugf("Got map: %+v", t)
	msg := &CommandMessage{}
	err = msg.SetCommand(t[types.CommandKey])
	if err != nil {
		log.Error(err)
		return nil, err
	}
	err = msg.SetArgs(t[types.ArgsKey])
	if err != nil {
		log.Error(err)
		return nil, err
	}
	return msg, nil
}
