package messages

import (
	"errors"

	erlang "github.com/okeuday/erlang_go/v2/erlang"
	log "github.com/sirupsen/logrus"

	"github.com/ut-proj/midiserver/pkg/erl/datatypes"
	"github.com/ut-proj/midiserver/pkg/types"
)

type CommandMessage struct {
	command types.CommandType
	args    map[string]string
}

func NewCommandMessage(t interface{}) (*CommandMessage, error) {
	tuple, ok := t.(erlang.OtpErlangTuple)
	if !ok {
		log.Debug("not tuple; checking to see if list of tuples ...")
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

func (cm *CommandMessage) Args() map[string]interface{} {
	return datatypes.MapStrsToInterfaces(cm.args)
}

func (cm *CommandMessage) SetCommand(cmd interface{}) error {
	cmdAtom, ok := cmd.(erlang.OtpErlangAtom)
	if !ok {
		return errors.New("could not cast command to atom")
	}
	cm.command = types.Command(types.CommandName(string(cmdAtom)))
	return nil
}

func (cm *CommandMessage) SetArgs(term interface{}) error {
	args, err := datatypes.TupleListToMap(term.(erlang.OtpErlangList))
	if err != nil {
		return err
	}
	cm.args = args
	return nil
}

func handleTuple(tuple erlang.OtpErlangTuple) (*CommandMessage, error) {
	log.Debug("handling tuple ...")
	tpl, err := datatypes.NewTupleFromTerm(tuple)
	if err != nil {
		log.Error(err)
		return nil, err
	}
	log.Debugf("Key: %+v (type %T)", tpl.Key(), tpl.Key())
	if tpl.Key() == types.CommandKey {
		msg := &CommandMessage{}
		err = msg.SetCommand(tpl.Value())
		if err != nil {
			log.Error(err)
			return nil, err
		}
		return msg, nil
	}
	return nil, nil
}

func handleTuples(tuples erlang.OtpErlangList) (*CommandMessage, error) {
	log.Debug("handling tuples ...")
	t, err := datatypes.TupleListToMap(tuples)
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
