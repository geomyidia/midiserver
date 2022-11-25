package midi

import (
	log "github.com/sirupsen/logrus"

	"github.com/ut-proj/midiserver/pkg/erl/datatypes"
)

// MIDI messages are expected to be in the format as produced by the
// midilib Erlang library (the `midimsg` module, in particular):
// * https://github.com/erlsci/midilib
//
// In summary, the following general forms need to be supported:
// * {midi, {batch, [{midi, {...}}]}}
// * {midi, {...}}
//
// The {...} parts of the messages have the following forms:
// * {command, Arg}
// * {command, [{ArgName, ArgValue}, ...]}

func HandleMessage(args interface{}) error {
	switch t := args.(type) {
	case *datatypes.List:
		log.Debugf("Got list: %+v", t.Elements())
	}
	return nil
}

// import (
// 	"errors"
// 	"fmt"
// 	"strconv"

// 	"github.com/google/uuid"
// 	erlang "github.com/okeuday/erlang_go/v2/erlang"
// 	log "github.com/sirupsen/logrus"

// 	"github.com/ut-proj/midiserver/pkg/erl/datatypes"
// 	"github.com/ut-proj/midiserver/pkg/types"
// )

// const MIDIKey = "midi"

// type MessageHandler struct {
// 	calls MIDICallGroup
// }

// type MIDICallGroup struct {
// 	id         string
// 	isParallel bool
// 	calls      []types.MidiCall
// }

// // This used to be part of NewMessageProcessor (now NewHandler):
// //
// // TODO: then check to see if a callBefore callback is set in the
// // options -- if so, pass the parsed datatypes interface (t) to it
// // and return that ... or whatever is best. Hrm, might be worth
// // defining an interface and letting projects implement that.
// // if t.HasKey(MIDIKey) {
// // 	mp.IsMidi = true
// // 	callGroup, err := NewMIDICallGroup(t)
// // 	if err != nil {
// // 		resp, err := NewResponse(types.EmptyResult, types.Err(err.Error()))
// // 		if err != nil {
// // 			return nil, err
// // 		}
// // 		resp.Send()
// // 		return nil, err
// // 	}
// // 	mp.midiCalls = callGroup
// // 	return mp, nil
// // }

// // This used to be part of the message processor (the Process method):
// //
// // if mp.cmdMsg != nil {
// // 	return types.Result(mp.cmdMsg.Command())
// // } else if mp.midiCalls.Length() != 0 {
// // 	return types.Result(types.MidiOp(types.MidiKey))
// // } else {
// // 	log.Error("unexpected message type")
// // 	return types.ContinueResult
// // }

// func NewMIDICallGroup(t interface{}) (*MIDICallGroup, error) {
// 	callGroup, err := Convert(t)
// 	if err != nil {
// 		log.Error(err)
// 		return nil, err
// 	}
// 	log.Debugf("converted calls: %+v", callGroup.calls)
// 	log.Debugf("parallel: %+v", callGroup.isParallel)
// 	return callGroup, nil
// }

// func (mcg *MIDICallGroup) Id() string {
// 	var id string
// 	if mcg != nil {
// 		id = mcg.id
// 	}
// 	return id
// }

// func (mcg *MIDICallGroup) IsParallel() bool {
// 	var parallel bool
// 	if mcg != nil {
// 		parallel = mcg.isParallel
// 	}
// 	return parallel
// }

// func (mcg *MIDICallGroup) Length() int {
// 	return len(mcg.calls)
// }

// func (mcg *MIDICallGroup) Calls() []types.MidiCall {
// 	var calls []types.MidiCall
// 	if mcg != nil {
// 		calls = mcg.calls
// 	}
// 	return calls
// }

// func Convert(term interface{}) (*MIDICallGroup, error) {
// 	var id string
// 	var parallel bool
// 	calls := []types.MidiCall{}
// 	data, err := datatypes.FromTerm(term)
// 	if err != nil {
// 		return nil, err
// 	}
// 	switch t := data.(type) {
// 	default:
// 		return nil, fmt.Errorf("could not convert %T", t)
// 	case *datatypes.List:
// 		for i, op := range t.Elements() {
// 			callGroup, err := Convert(op)
// 			if err != nil {
// 				log.Error(err)
// 				return nil, err
// 			}
// 			var updatedCall []types.MidiCall
// 			for _, op := range callGroup.calls {
// 				op.Id = i + 1
// 				updatedCall = append(updatedCall, op)
// 			}
// 			calls = append(calls, updatedCall...)
// 		}
// 		return &MIDICallGroup{calls: calls}, nil
// 	case *datatypes.Tuple:
// 		log.Debugf("tuple data: %+v", t)
// 		var tpl *datatypes.Tuple
// 		var ok bool
// 		// If the tuple is a MIDI command, then it's going to be
// 		// a tuple of tuples:
// 		if t.Key().(*datatypes.Atom).Value() == MIDIKey {
// 			tpl, ok = t.Value().(*datatypes.Tuple)
// 			if !ok {
// 				return nil, datatypes.ErrCastingTuple
// 			}
// 			log.Debugf("MIDIKey tuple: %+v", tpl)
// 		}
// 		key := tpl.Key().(*datatypes.Atom).Value()
// 		if key == types.MidiBatchKey {
// 			batchCallGroup, err := ConvertBatch(tpl.Value())
// 			if err != nil {
// 				log.Error(err)
// 				return nil, err
// 			}
// 			id = batchCallGroup.id
// 			parallel = batchCallGroup.isParallel
// 			log.Debug("batch parallel: ", parallel)
// 			calls = append(calls, batchCallGroup.calls...)
// 		} else {
// 			args, err := ConvertArg(key, tpl.Value())
// 			if err != nil {
// 				log.Error(err)
// 				return nil, err
// 			}
// 			call := types.MidiCall{Op: types.MidiOpType(key), Args: args}
// 			calls = append(calls, call)
// 		}
// 		log.Debug("parallel: ", parallel)
// 		return &MIDICallGroup{
// 			id:         id,
// 			isParallel: parallel,
// 			calls:      calls,
// 		}, nil
// 	}
// }

// func ConvertBatch(term interface{}) (*MIDICallGroup, error) {
// 	var id string
// 	var parallel bool
// 	data, err := datatypes.FromTerm(term)
// 	if err != nil {
// 		return nil, err
// 	}
// 	list, ok := data.(*datatypes.List)
// 	if !ok {
// 		return nil, datatypes.ErrCastingList
// 	}
// 	log.Debugf("converting batch: %+v", list)

// 	batchMap, err := datatypes.TupleListToMap(list)
// 	if err != nil {
// 		return nil, err
// 	}
// 	log.Debugf("batch map data: %+v", batchMap)

// 	// Process the Batch ID
// 	rawId := batchMap[types.MidiIdKey]
// 	uuid4, err := uuid.FromBytes([]byte(rawId))
// 	if err != nil {
// 		return nil, err
// 	}
// 	id = uuid4.String()

// 	// Get parallel flag
// 	rawParallel := batchMap[types.MidiParallelKey]
// 	if rawParallel != "" {
// 		parallel, err = strconv.ParseBool(rawParallel)
// 		log.Debug("parallel? ", parallel)
// 		if err != nil {
// 			return nil, err
// 		}
// 	}

// 	// Process the Batch Messages
// 	batch, err := Convert(batchMap[types.MidiMessagesKey])
// 	if err != nil {
// 		return nil, err
// 	}
// 	batch.id = id
// 	batch.isParallel = parallel
// 	log.Debug("parallel: ", parallel)
// 	log.Debug("set parallel: ", batch.isParallel)
// 	return batch, nil
// }

// func ConvertArg(k string, v interface{}) (*types.MidiArgs, error) {
// 	log.Debug("converting args ...")
// 	args := &types.MidiArgs{}
// 	switch k {
// 	case types.MidiChannelKey:
// 		args.Channel = v.(uint8)
// 	case types.MidiDeviceKey:
// 		args.Device = v.(uint8)
// 	case types.MidiNoteOffKey:
// 		args.NoteOff = v.(uint8)
// 	case types.MidiNoteOnKey:
// 		list := v.(erlang.OtpErlangList)
// 		noteOnData, err := datatypes.TupleListToMap(list)
// 		if err != nil {
// 			log.Error(err)
// 			return nil, err
// 		}
// 		log.Debugf("noteOnData: %+v", noteOnData)
// 		noteOn := datatypes.MapStrsToInterfaces(noteOnData)
// 		args.NoteOn = types.MidiNoteOn{
// 			Pitch:    noteOn[types.MidiPitchKey].(uint8),
// 			Velocity: noteOn[types.MidiVelocityKey].(uint8),
// 		}
// 	case types.MidiProgramChangeKey:
// 		args.Program = v.(uint8)
// 	case types.MidiBankSelectMSBKey, types.MidiBankSelectLSBKey:
// 		args.CC = types.MidiCC{
// 			Value: v.(uint8),
// 		}
// 	case types.MidiCCKey:
// 		list := v.(erlang.OtpErlangList)
// 		ccData, err := datatypes.TupleListToMap(list)
// 		if err != nil {
// 			log.Error(err)
// 			return nil, err
// 		}
// 		cc := datatypes.MapStrsToInterfaces(ccData)
// 		args.CC = types.MidiCC{
// 			Controller: cc[types.MidiControllerKey].(uint8),
// 			Value:      cc[types.MidiValueKey].(uint8),
// 		}
// 	case types.MidiRealtimeKey:
// 		atomRealtime, ok := v.(erlang.OtpErlangAtom)
// 		if !ok {
// 			return nil, errors.New("couldn't convert 'realtime'")
// 		}
// 		args.Realtime = types.MidiRTType(string(atomRealtime))
// 	}
// 	return args, nil
// }
