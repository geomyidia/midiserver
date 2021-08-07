package messages

import (
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
	calls, err := Convert(t)
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

func Convert(term interface{}) ([]types.MidiCall, error) {
	emptyCalls := []types.MidiCall{}
	calls := []types.MidiCall{}
	switch t := term.(type) {
	default:
		return emptyCalls, fmt.Errorf("could not convert %T", t)
	case erlang.OtpErlangList:
		ops, ok := term.(erlang.OtpErlangList)
		fmt.Printf("%+v\n", ops)
		if !ok {
			return emptyCalls, fmt.Errorf("could not convert %T", t)
		}
		for _, op := range ops.Value {
			call, err := Convert(op)
			if err != nil {
				log.Error(err)
				return emptyCalls, err
			}
			calls = append(calls, call...)
		}
		return calls, nil
	case erlang.OtpErlangTuple:
		key, val, err := datatypes.Tuple(t)
		if err != nil {
			log.Error(err)
			return emptyCalls, err
		}
		if key == types.MidiKey {
			key, val, err = datatypes.Tuple(val)
		}
		if err != nil {
			log.Error(err)
			return emptyCalls, err
		}
		if key == types.MidiBatchKey {
			batchCalls, err := Convert(val)
			if err != nil {
				log.Error(err)
				return emptyCalls, err
			}
			calls = append(calls, batchCalls...)
		} else {
			args, err := ConvertArg(key, val)
			if err != nil {
				log.Error(err)
				return emptyCalls, err
			}
			call := types.MidiCall{Op: types.MidiOpType(key), Args: args}
			calls = append(calls, call)
		}
		return calls, nil
	}
}
