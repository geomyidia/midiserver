package messages

import (
	"errors"

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
		args, err := convertArgs(val)
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

func convertArgs(t interface{}) (*types.MidiArgs, error) {
	log.Debug("Converting args ...")
	args := &types.MidiArgs{}
	propList, ok := t.(erlang.OtpErlangList)
	if !ok {
		return args, errors.New("couldn't parse batches")
	}
	data, err := datatypes.PropListToMap(propList)
	if err != nil {
		log.Error(err)
		return nil, err
	}
	for k, v := range data {
		switch k {
		case "device":
			args.Device = v.(uint8)
		case "tempo":
			args.Tempo = v.(uint8)
		case "note_off":
			args.NoteOff = v.(uint8)
		case "meter":
			list := v.(erlang.OtpErlangList)
			meter, err := datatypes.PropListToMap(list)
			if err != nil {
				log.Error(err)
				continue
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
				continue
			}
			args.NoteOn = types.MidiNoteOn{
				Pitch:    noteOn["pitch"].(uint8),
				Velocity: noteOn["velocity"].(uint8),
			}
		}
	}
	return args, nil
}
