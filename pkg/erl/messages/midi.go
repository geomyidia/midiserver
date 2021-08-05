package messages

import (
	erlang "github.com/okeuday/erlang_go/v2/erlang"
	log "github.com/sirupsen/logrus"

	"github.com/geomyidia/midiserver/pkg/erl/datatypes"
	"github.com/geomyidia/midiserver/pkg/types"
)

type MidiMessage struct {
	op   types.MidiOpType
	data types.PropList
	ops  types.MidiOps
	args *types.MidiArgs
}

func NewMidiMessage(t interface{}) (*MidiMessage, error) {
	msg := &MidiMessage{}
	op, data, err := MidiOp(t)
	if err != nil {
		return nil, err
	}
	msg.op = op
	msg.handleData(data)
	log.Tracef("Got args: %+v", msg.args)
	return msg, nil
}

func MidiOp(term interface{}) (types.MidiOpType, types.PropList, error) {
	_, opTerm, err := datatypes.Tuple(term)
	if err != nil {
		return "", nil, err
	}
	key, val, err := datatypes.Tuple(opTerm)
	if err != nil {
		return "", nil, err
	}
	data, err := datatypes.PropListToMap(val.(erlang.OtpErlangList))
	if err != nil {
		return "", nil, err
	}
	return types.MidiOp(key), data, nil
}

func (mm *MidiMessage) Args() *types.MidiArgs {
	return mm.args
}

func (mm *MidiMessage) Op() types.MidiOpType {
	return mm.op
}

func (mm *MidiMessage) Data() types.PropList {
	return mm.data
}

func (mm *MidiMessage) handleData(data types.PropList) {
	log.Debug("Handling op data ...")
	switch mm.op {
	case types.MidiBatch():
		mm.handleBatch(data)
	}
	mm.data = data
}

func (mm *MidiMessage) handleBatch(data types.PropList) {
	log.Debug("Handling batch op ...")
	args := &types.MidiArgs{}
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
	mm.args = args
}
