package messages

import (
	log "github.com/sirupsen/logrus"

	"github.com/ut-proj/midiserver/pkg/erl/datatypes"
)

type Message struct {
	messageType *datatypes.Atom
	name        *datatypes.Atom
	args        *datatypes.List
}

func NewFromBytes(data []byte) (*Message, error) {
	t, err := datatypes.FromBytes(data)
	if err != nil {
		log.Error(err)
		return nil, err
	}
	return New(t)
}

func NewFromTerm(term interface{}) (*Message, error) {
	t, err := datatypes.FromTerm(term)
	if err != nil {
		log.Error(err)
		return nil, err
	}
	return New(t)
}

func New(t interface{}) (*Message, error) {
	var ok bool
	msgTuple, err := messageTuple(t)
	if err != nil {
		return nil, err
	}
	log.Debugf("Got message tuple: %+v", msgTuple)
	msgType, ok := msgTuple.Key().(*datatypes.Atom)
	if !ok {
		log.Error(ErrMsgAtomFormat)
		return nil, ErrMsgAtomFormat
	}
	log.Debugf("Got message type %v", msgType)

	args := datatypes.NewList([]interface{}{})
	name, ok := msgTuple.Value().(*datatypes.Atom)
	if ok {
		// This is the case for a simple message with just a name and no args
		log.Debugf("Got message name %s", name)
	} else {
		// This is the case for more highly-structured messages
		name, args, err = messageNameArgs(msgTuple)
		if err != nil {
			return nil, err
		}
	}
	// If this is a MIDI message
	return &Message{
		messageType: msgType,
		name:        name,
		args:        args,
	}, nil
}

func NewCommandFromName(name string) *Message {
	return &Message{
		messageType: datatypes.NewAtom("command"),
		name:        datatypes.NewAtom(name),
	}
}

func (m *Message) Type() string {
	return m.messageType.Value()
}

func (m *Message) Name() string {
	return m.name.Value()
}

func (m *Message) Args() []interface{} {
	return m.args.Elements()
}

// Private functions

func messageTuple(t interface{}) (*datatypes.Tuple, error) {
	var msgTuple *datatypes.Tuple
	log.Tracef("Got Go/Erlang ports data: %+v", t)
	parts, ok := t.(*datatypes.List)
	if ok {
		if parts.Len() > 2 {
			log.Error(ErrMsgListFormat)
			return nil, ErrMsgListFormat
		}
		msgTuple, ok = parts.Nth(0).(*datatypes.Tuple)
		if !ok {
			log.Error(ErrMsgTupleFormat)
			return nil, ErrMsgTupleFormat
		}
	} else {
		msgTuple, ok = t.(*datatypes.Tuple)
		if !ok {
			log.Error(ErrMsgTupleFormat)
			return nil, ErrMsgTupleFormat
		}
	}
	return msgTuple, nil
}

func messageNameArgs(msgTuple *datatypes.Tuple) (*datatypes.Atom, *datatypes.List, error) {
	var msgData []interface{}
	var name *datatypes.Atom
	var ok bool
	msgVal := msgTuple.Value()
	x, ok := msgVal.(*datatypes.List)
	if ok {
		msgData = x.Elements()
	} else {
		x, ok := msgVal.(*datatypes.Tuple)
		if !ok {
			log.Error(ErrMsgValueFormat)
			return nil, nil, ErrMsgValueFormat
		}
		msgData = x.Elements()
	}
	if len(msgData) == 0 {
		log.Error(ErrMsgValueFormat)
		return nil, nil, ErrMsgValueFormat
	}
	name, ok = msgData[0].(*datatypes.Atom)
	if !ok {
		log.Error(ErrMsgNameFormat)
		return nil, nil, ErrMsgNameFormat
	}
	log.Debugf("Got message name %s", name)
	args := datatypes.NewList(msgData[1:])
	return name, args, nil
}
