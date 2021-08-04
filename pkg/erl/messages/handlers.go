package messages

import (
	erlang "github.com/okeuday/erlang_go/v2/erlang"
	log "github.com/sirupsen/logrus"

	"github.com/geomyidia/midiserver/pkg/erl/proplists"
	"github.com/geomyidia/midiserver/pkg/types"
)

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
