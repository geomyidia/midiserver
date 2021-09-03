package midi

import (
	"context"
	"errors"
	"fmt"

	log "github.com/sirupsen/logrus"
	"gitlab.com/gomidi/midi"
	"gitlab.com/gomidi/midi/midimessage/channel"
	"gitlab.com/gomidi/midi/writer"
	"gitlab.com/gomidi/rtmididrv"

	"github.com/ut-proj/midiserver/pkg/types"
)

type System struct {
	Driver          *rtmididrv.Driver
	DevicesIn       []midi.In
	DevicesOut      []midi.Out
	DeviceIn        midi.In
	DeviceOut       midi.Out
	Writer          *writer.Writer
	DeviceOutOpened bool
	ChannelSet      bool
}

func NewSystem() *System {
	log.Info("creating MIDI system ...")
	drv, err := rtmididrv.New()
	if err != nil {
		log.Fatal(err)
	}

	ins, err := drv.Ins()
	if err != nil {
		log.Fatal(err)
	}

	outs, err := drv.Outs()
	if err != nil {
		log.Fatal(err)
	}

	return &System{
		Driver:          drv,
		DevicesIn:       ins,
		DevicesOut:      outs,
		DeviceOutOpened: false,
		ChannelSet:      false,
	}
}

func (s *System) Shutdown() {
	log.Info("shutting down MIDI system ...")
	s.Driver.Close()
	s.DeviceOutOpened = false
}

func (s *System) SetDevice(deviceId uint8) error {
	if s.DeviceOutOpened {
		return nil
	}
	log.Info("setting device ...")
	s.DeviceOut = s.DevicesOut[deviceId]
	err := s.DeviceOut.Open()
	if err != nil {
		log.Fatal(err)
		return err
	}
	s.Writer = writer.New(s.DeviceOut)
	s.DeviceOutOpened = true
	return nil
}

func (s *System) SetChannel(channelId uint8) error {
	log.Info("setting channel ...")
	s.Writer.SetChannel(channelId)
	s.ChannelSet = true
	return nil
}

func (s *System) GetChannel() channel.Channel {
	return channel.Channel(s.Writer.Channel())
}

func (s *System) Dispatch(ctx context.Context, calls []types.MidiCall,
	isParallel bool, flags *types.Flags) {
	var err error
	log.Trace("dispatching MIDI operation ...")
	log.Debugf("got MIDI calls: %v", calls)

	for _, call := range calls {
		log.Debugf("making MIDI call '%s' ...", call.Op)
		log.Debugf("with (id, args): (%d, %+v) ...", call.Id, call.Args)
		// XXX This isn't safe ... what's the right way to do this in gomidi?
		// if isParallel {
		// 	log.Trace("making calling in parallel ...")
		// 	go s.CallMidi(call)
		// } else {
		// 	log.Trace("making calling in series ...")
		// 	err = s.CallMidi(call)
		// }
		err = s.CallMidi(call)
		if err != nil {
			log.Error(err)
		}
	}
}

func (s *System) CallMidi(call types.MidiCall) error {
	switch call.Op {
	case types.MidiDeviceType():
		return s.SetDevice(call.Args.Device)
	case types.MidiChannelType():
		return s.SetChannel(call.Args.Channel)
	}
	if !s.DeviceOutOpened {
		return errors.New("can't send command when device not opened")
	}
	if !s.ChannelSet {
		return fmt.Errorf("message of type %s require the channel to be set",
			call.Op)
	}
	switch call.Op {
	case types.MidiNoteOnType():
		log.Tracef("calling NoteOn with values: %+v", call.Args.NoteOn)
		err := writer.NoteOn(s.Writer, call.Args.NoteOn.Pitch, call.Args.NoteOn.Velocity)
		if err != nil {
			return err
		}
		return nil
	case types.MidiNoteOffType():
		err := writer.NoteOff(s.Writer, call.Args.NoteOff)
		if err != nil {
			return err
		}
		return nil
	case types.MidiCCType():
		err := writer.ControlChange(s.Writer, call.Args.CC.Controller, call.Args.CC.Value)
		if err != nil {
			return err
		}
		return nil
	default:
		log.Errorf("no handler for operation '%s'", call.Op)
		return nil
	}
}
