package messages

import (
	"errors"

	erlang "github.com/okeuday/erlang_go/v2/erlang"

	"github.com/ut-proj/midiserver/pkg/erl/datatypes"
	"github.com/ut-proj/midiserver/pkg/types"
)

type CommandMessage struct {
	command types.CommandType
	args    types.PropList
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

func (cm *CommandMessage) Args() types.PropList {
	return cm.args
}

func (cm *CommandMessage) SetCommand(cmdIf interface{}) error {
	cmdAtom, ok := cmdIf.(erlang.OtpErlangAtom)
	if !ok {
		return errors.New("could not cast command to atom")
	}
	cm.command = types.Command(types.CommandName(string(cmdAtom)))
	return nil
}

func (cm *CommandMessage) SetArgs(argsIf interface{}) error {
	args, err := datatypes.PropListToMap(argsIf.(erlang.OtpErlangList))
	if err != nil {
		return err
	}
	cm.args = args
	return nil
}
