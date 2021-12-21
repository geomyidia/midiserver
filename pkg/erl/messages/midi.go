package messages

import (
	"errors"
	"fmt"
	"strconv"

	"github.com/google/uuid"
	erlang "github.com/okeuday/erlang_go/v2/erlang"
	log "github.com/sirupsen/logrus"

	"github.com/ut-proj/midiserver/pkg/erl/datatypes"
	"github.com/ut-proj/midiserver/pkg/types"
)

type MidiCallGroup struct {
	id         string
	isParallel bool
	calls      []types.MidiCall
}

func NewMidiCallGroup(t interface{}) (*MidiCallGroup, error) {
	callGroup, err := Convert(t)
	if err != nil {
		log.Error(err)
		return nil, err
	}
	log.Debugf("converted calls: %+v", callGroup.calls)
	log.Debugf("parallel: %+v", callGroup.isParallel)
	return callGroup, nil
}

func (mcg *MidiCallGroup) Id() string {
	var id string
	if mcg != nil {
		id = mcg.id
	}
	return id
}

func (mcg *MidiCallGroup) IsParallel() bool {
	var parallel bool
	if mcg != nil {
		parallel = mcg.isParallel
	}
	return parallel
}

func (mcg *MidiCallGroup) Length() int {
	return len(mcg.calls)
}

func (mcg *MidiCallGroup) Calls() []types.MidiCall {
	var calls []types.MidiCall
	if mcg != nil {
		calls = mcg.calls
	}
	return calls
}

func Convert(term interface{}) (*MidiCallGroup, error) {
	var id string
	var parallel bool
	calls := []types.MidiCall{}
	switch t := term.(type) {
	default:
		return nil, fmt.Errorf("could not convert %T", t)
	case erlang.OtpErlangList:
		ops, ok := term.(erlang.OtpErlangList)
		if !ok {
			return nil, fmt.Errorf("could not convert %T", t)
		}
		for idx, op := range ops.Value {
			callGroup, err := Convert(op)
			if err != nil {
				log.Error(err)
				return nil, err
			}
			var updatedCall []types.MidiCall
			for _, op := range callGroup.calls {
				op.Id = idx + 1
				updatedCall = append(updatedCall, op)
			}
			calls = append(calls, updatedCall...)
		}
		return &MidiCallGroup{calls: calls}, nil
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
		if key == types.MidiBatchKey {
			batchCallGroup, err := ConvertBatch(val)
			if err != nil {
				log.Error(err)
				return nil, err
			}
			id = batchCallGroup.id
			parallel = batchCallGroup.isParallel
			log.Debug("batch parallel: ", parallel)
			calls = append(calls, batchCallGroup.calls...)
		} else {
			args, err := ConvertArg(key, val)
			if err != nil {
				log.Error(err)
				return nil, err
			}
			call := types.MidiCall{Op: types.MidiOpType(key), Args: args}
			calls = append(calls, call)
		}
		log.Debug("parallel: ", parallel)
		return &MidiCallGroup{
			id:         id,
			isParallel: parallel,
			calls:      calls,
		}, nil
	}
}

func ConvertBatch(term interface{}) (*MidiCallGroup, error) {
	var id string
	var parallel bool
	list, ok := term.(erlang.OtpErlangList)
	log.Debugf("converting batch: %+v", list)
	if !ok {
		return nil, errors.New("couldn't convert batch")
	}
	batchMap, err := datatypes.PropListToMap(list)
	if err != nil {
		return nil, err
	}
	// Process the Batch ID
	rawId := batchMap[types.MidiIdKey]
	binId, ok := rawId.(erlang.OtpErlangBinary)
	if !ok {
		return nil, errors.New("couldn't convert batch id")
	}
	rawParallel := batchMap[types.MidiParallelKey]
	if rawParallel != nil {
		atomParallel, ok := rawParallel.(erlang.OtpErlangAtom)
		if !ok {
			return nil, errors.New("couldn't convert 'parallel?'")
		}
		parallel, err = strconv.ParseBool(string(atomParallel))
		log.Debug("parallel? ", parallel)
		if err != nil {
			return nil, err
		}
	}
	uuid4, err := uuid.FromBytes(binId.Value)
	if err != nil {
		return nil, err
	}
	id = uuid4.String()
	// Process the Batch Messages
	batch, err := Convert(batchMap[types.MidiMessagesKey])

	if err != nil {
		return nil, err
	}
	batch.id = id
	batch.isParallel = parallel
	log.Debug("parallel: ", parallel)
	log.Debug("set parallel: ", batch.isParallel)
	return batch, nil
}

func ConvertArg(k string, v interface{}) (*types.MidiArgs, error) {
	log.Debug("converting args ...")
	args := &types.MidiArgs{}
	switch k {
	case types.MidiChannelKey:
		args.Channel = v.(uint8)
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
			Pitch:    noteOn[types.MidiPitchKey].(uint8),
			Velocity: noteOn[types.MidiVelocityKey].(uint8),
		}
	case types.MidiProgramChangeKey:
		args.Program = v.(uint8)
	case types.MidiBankSelectMSBKey, types.MidiBankSelectLSBKey:
		args.CC = types.MidiCC{
			Value: v.(uint8),
		}
	case types.MidiCCKey:
		list := v.(erlang.OtpErlangList)
		cc, err := datatypes.PropListToMap(list)
		if err != nil {
			log.Error(err)
			return nil, err
		}
		args.CC = types.MidiCC{
			Controller: cc[types.MidiControllerKey].(uint8),
			Value:      cc[types.MidiValueKey].(uint8),
		}
	case types.MidiRealtimeKey:
		atomRealtime, ok := v.(erlang.OtpErlangAtom)
		if !ok {
			return nil, errors.New("couldn't convert 'realtime'")
		}
		args.Realtime = types.MidiRTType(string(atomRealtime))
	}
	return args, nil
}
