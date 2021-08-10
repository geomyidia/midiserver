package messages

import (
	"errors"
	"fmt"

	"github.com/google/uuid"
	erlang "github.com/okeuday/erlang_go/v2/erlang"
	log "github.com/sirupsen/logrus"

	"github.com/geomyidia/midiserver/pkg/erl/datatypes"
	"github.com/geomyidia/midiserver/pkg/types"
)

type MidiCallGroup struct {
	id    string
	calls []types.MidiCall
}

func NewMidiCalls(t interface{}) (*MidiCallGroup, error) {
	id, calls, err := Convert(t)
	if err != nil {
		log.Error(err)
		return nil, err
	}
	return &MidiCallGroup{id: id, calls: calls}, nil
}

func (mcg *MidiCallGroup) Id() string {
	return mcg.id
}

func (mcg *MidiCallGroup) Length() int {
	return len(mcg.calls)
}

func (mcg *MidiCallGroup) Calls() []types.MidiCall {
	return mcg.calls
}

func ConvertArg(k string, v interface{}) (*types.MidiArgs, error) {
	args := &types.MidiArgs{}
	switch k {
	case types.MidiDeviceKey:
		args.Device = v.(uint8)
	case types.MidiNoteOffKey:
		args.NoteOff = v.(uint8)
	case types.MidiNoteOnKey:
		list := v.(erlang.OtpErlangList)
		noteOn, err := datatypes.PropListToMap(list)
		if err != nil {
			log.Error(err)
			return nil, err
		}
		args.NoteOn = types.MidiNoteOn{
			Pitch:    noteOn["pitch"].(uint8),
			Velocity: noteOn["velocity"].(uint8),
		}
	}
	return args, nil
}

func Convert(term interface{}) (string, []types.MidiCall, error) {
	var id string
	calls := []types.MidiCall{}
	switch t := term.(type) {
	default:
		return "", nil, fmt.Errorf("could not convert %T", t)
	case erlang.OtpErlangList:
		ops, ok := term.(erlang.OtpErlangList)
		fmt.Printf("%+v\n", ops)
		if !ok {
			return "", nil, fmt.Errorf("could not convert %T", t)
		}
		for idx, op := range ops.Value {
			_, call, err := Convert(op)
			if err != nil {
				log.Error(err)
				return "", nil, err
			}
			var updatedCall []types.MidiCall
			for _, op := range call {
				op.Id = idx + 1
				updatedCall = append(updatedCall, op)
			}
			calls = append(calls, updatedCall...)
		}
		return "", calls, nil
	case erlang.OtpErlangTuple:
		key, val, err := datatypes.Tuple(t)
		if err != nil {
			log.Error(err)
			return "", nil, err
		}
		if key == types.MidiKey {
			key, val, err = datatypes.Tuple(val)
		}
		if err != nil {
			log.Error(err)
			return "", nil, err
		}
		if key == types.MidiBatchKey {
			var batchCalls []types.MidiCall
			id, batchCalls, err = ConvertBatch(val)
			if err != nil {
				log.Error(err)
				return "", nil, err
			}
			calls = append(calls, batchCalls...)
		} else {
			args, err := ConvertArg(key, val)
			if err != nil {
				log.Error(err)
				return "", nil, err
			}
			call := types.MidiCall{Op: types.MidiOpType(key), Args: args}
			calls = append(calls, call)
		}
		return id, calls, nil
	}
}

func ConvertBatch(term interface{}) (string, []types.MidiCall, error) {
	var id string
	list, ok := term.(erlang.OtpErlangList)
	if !ok {
		return "", nil, errors.New("couldn't convert batch")
	}
	batchMap, err := datatypes.PropListToMap(list)
	if err != nil {
		return "", nil, err
	}
	// Process the Batch ID
	rawId := batchMap[types.MidiIdKey]
	binId, ok := rawId.(erlang.OtpErlangBinary)
	if !ok {
		return "", nil, errors.New("couldn't convert batch id")
	}
	uuid4, err := uuid.FromBytes(binId.Value)
	if err != nil {
		return "", nil, err
	}
	id = uuid4.String()
	// Process the Batch Messages
	_, batch, err := Convert(batchMap[types.MidiMessagesKey])
	if err != nil {
		return "", nil, err
	}
	return id, batch, nil
}
