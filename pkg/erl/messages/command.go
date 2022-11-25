package messages

import (
	log "github.com/sirupsen/logrus"

	"github.com/ut-proj/midiserver/pkg/erl/datatypes"
)

type Command struct {
	command *datatypes.Atom
	args    *datatypes.List
}

func NewCommand(term interface{}) (*Command, error) {
	var ok bool
	var cmdTuple *datatypes.Tuple
	t, err := datatypes.FromTerm(term)
	if err != nil {
		log.Error(err)
		return nil, err
	}
	log.Tracef("got Go/Erlang ports data: %+v", t)
	cmdList, ok := t.(*datatypes.List)
	if ok {
		if cmdList.Len() > 2 {
			log.Error(ErrCmdListFormat)
			return nil, ErrCmdListFormat
		}
		cmdTuple, ok = cmdList.Nth(0).(*datatypes.Tuple)
		if !ok {
			log.Error(ErrCmdTupleFormat)
			return nil, ErrCmdTupleFormat
		}
	} else {
		cmdTuple, ok = t.(*datatypes.Tuple)
		if !ok {
			log.Error(ErrCmdTupleFormat)
			return nil, ErrCmdTupleFormat
		}
	}

	cmd, ok := cmdTuple.Key().(*datatypes.Atom)
	if !ok {
		log.Error(ErrCmdAtomFormat)
		return nil, ErrCmdAtomFormat
	}

	args := new(datatypes.List)
	arg, ok := cmdTuple.Value().(*datatypes.Atom)
	if ok {
		args.Append(arg)
	} else {
		args, ok = cmdTuple.Value().(*datatypes.List)
		if !ok {
			return nil, ErrCmdValueFormat
		}
	}
	return &Command{
		command: cmd,
		args:    args,
	}, nil
}

func (cm *Command) Name() string {
	return cm.command.Value()
}

func (cm *Command) Args() []interface{} {
	return cm.args.Elements()
}
