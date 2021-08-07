package messages

import (
	"errors"
	"fmt"

	erlang "github.com/okeuday/erlang_go/v2/erlang"
	log "github.com/sirupsen/logrus"

	"github.com/geomyidia/midiserver/pkg/erl/datatypes"
	"github.com/geomyidia/midiserver/pkg/types"
)

type MidiCalls struct {
	calls []types.MidiCall
}

func NewMidiCalls(t interface{}) (*MidiCalls, error) {
	calls, err := getCalls(t)
	if err != nil {
		log.Error(err)
		return nil, err
	}
	return &MidiCalls{calls: calls}, nil
}

func (mc *MidiCalls) Length() int {
	return len(mc.calls)
}

func (mc *MidiCalls) Calls() []types.MidiCall {
	return mc.calls
}

func getCalls(t interface{}) ([]types.MidiCall, error) {
	var calls []types.MidiCall
	key, val, err := datatypes.Tuple(t)
	if err != nil {
		log.Error(err)
		return calls, err
	}
	if key == types.MidiKey {
		return getCalls(val)
	}
	if key == types.MidiBatchCall {
		batches, ok := val.(erlang.OtpErlangList)
		if !ok {
			return calls, errors.New("couldn't parse batches")
		}
		log.Tracef("Batches: %+v", batches)
		for _, op := range batches.Value {
			subCalls, err := getCalls(op)
			log.Tracef("Sub-calls: %+v", subCalls)
			if err != nil {
				log.Error(err)
				return calls, err
			}
			calls = append(calls, subCalls...)
		}
	} else {
		op := types.MidiOpType(key)
		log.Tracef("Op: %v", op)
		args, err := ConvertArgs(val)
		if err != nil {
			log.Error(err)
			return calls, err
		}
		log.Tracef("Args: %+v", args)
		calls = append(calls, types.MidiCall{Op: op, Args: args})
		log.Tracef("Calls: %+v", calls)
		return calls, nil
	}
	return calls, nil
}

func ConvertArgs(t interface{}) (*types.MidiArgs, error) {
	log.Debug("Converting args ...")
	var args *types.MidiArgs
	nilArgs := &types.MidiArgs{}
	propList, ok := t.(erlang.OtpErlangList)
	if !ok {
		return nilArgs, errors.New("couldn't parse batches")
	}
	data, err := datatypes.PropListToMap(propList)
	if err != nil {
		log.Error(err)
		return nil, err
	}
	for k, v := range data {
		args, err = ConvertArg(k, v)
		if err != nil {
			continue
		}
	}
	return args, nil
}

func ConvertArg(k string, v interface{}) (*types.MidiArgs, error) {
	args := &types.MidiArgs{}
	switch k {
	case "device":
		args.Device = v.(uint8)
	case "tempo_bpm":
		args.Tempo = v.(uint8)
	case "note_off":
		args.NoteOff = v.(uint8)
	case "meter":
		list := v.(erlang.OtpErlangList)
		meter, err := datatypes.PropListToMap(list)
		if err != nil {
			log.Error(err)
			return nil, err
		}
		args.Meter = types.MidiMeter{
			Numerator:   meter["numerator"].(uint8),
			Denominator: meter["denominator"].(uint8),
		}
	case "note_on":
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

func Convert(term interface{}) (*types.MidiCall, error) {
	switch t := term.(type) {
	default:
		return nil, fmt.Errorf("Could not convert %T", t)
	case erlang.OtpErlangTuple:
		key, val, err := datatypes.Tuple(t)
		if err != nil {
			log.Error(err)
			return nil, err
		}
		if key == types.MidiKey {
			key, val, err = datatypes.Tuple(val)
		}
		if err != nil {
			log.Error(err)
			return nil, err
		}
		args, err := ConvertArg(key, val)
		if err != nil {
			log.Error(err)
			return nil, err
		}
		return &types.MidiCall{Op: types.MidiOpType(key), Args: args}, nil
	}
}
