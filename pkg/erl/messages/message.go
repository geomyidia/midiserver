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

func New(term interface{}) (*Message, error) {
	var ok bool
	var msgTuple *datatypes.Tuple
	t, err := datatypes.FromTerm(term)
	if err != nil {
		log.Error(err)
		return nil, err
	}
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
	log.Debugf("Got message tuple: %+v", msgTuple)
	msgType, ok := msgTuple.Key().(*datatypes.Atom)
	if !ok {
		log.Error(ErrMsgAtomFormat)
		return nil, ErrMsgAtomFormat
	}
	log.Debugf("Got message type %v", msgType)

	var name *datatypes.Atom
	args := datatypes.NewList([]interface{}{})
	name, ok = msgTuple.Value().(*datatypes.Atom)
	if ok {
		log.Debugf("Got message name %s", name)
	} else {
		msgData, ok := msgTuple.Value().(*datatypes.List)
		if !ok {
			log.Error(ErrMsgValueFormat)
			return nil, ErrMsgValueFormat
		}
		if len(msgData.Elements()) == 0 {
			log.Error(ErrMsgValueFormat)
			return nil, ErrMsgValueFormat
		}
		name, ok = msgData.Elements()[0].(*datatypes.Atom)
		if !ok {
			log.Error(ErrMsgNameFormat)
			return nil, ErrMsgNameFormat
		}
		log.Debugf("Got message name %s", name)
		args = datatypes.NewList(msgData.Elements()[1:])
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
