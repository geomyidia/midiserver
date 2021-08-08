package midi

import (
	"context"
	"errors"

	log "github.com/sirupsen/logrus"
	"gitlab.com/gomidi/midi"
	"gitlab.com/gomidi/midi/writer"
	"gitlab.com/gomidi/rtmididrv"

	"github.com/geomyidia/midiserver/pkg/types"
)

type System struct {
	Driver          *rtmididrv.Driver
	DevicesIn       []midi.In
	DevicesOut      []midi.Out
	DeviceIn        midi.In
	DeviceOut       midi.Out
	DeviceOutOpened bool
	Writer          *writer.Writer
}

func NewSystem() *System {
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
		Driver:     drv,
		DevicesIn:  ins,
		DevicesOut: outs,
	}
}

func (s *System) Shutdown() {
	log.Info("shutting down MIDI system ...")
	s.Driver.Close()
	s.DeviceOutOpened = false
}

func (s *System) SetDevice(deviceId uint8) error {
	log.Trace("setting device ...")
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

func (s *System) SetChannel(channelId uint8) {
	log.Trace("setting channel ...")
	s.Writer.SetChannel(channelId)
}

func (s *System) GetChannel(channelId uint8) uint8 {
	return s.Writer.Channel()
}

func (s *System) Dispatch(ctx context.Context, calls []types.MidiCall, flags *types.Flags) {
	log.Debug("dispatching MIDI operation ...")
	log.Tracef("got MIDI calls: %v", calls)
	for _, call := range calls {
		log.Debugf("making MIDI call %v ...", call)
		err := s.CallMidi(call)
		if err != nil {
			log.Error(err)
		}
	}
}

func (s *System) CallMidi(call types.MidiCall) error {
	switch call.Op {
	case types.MidiDeviceType():
		s.SetDevice(call.Args.Device)
		return nil
	case types.MidiChannelType():
		s.SetChannel(call.Args.Channel)
		return nil
	case types.MidiMeterType():
		println("tbd")
		return nil
	case types.MidiTempoType():
		println("tbd")
		return nil
	case types.MidiNoteOnType():
		log.Tracef("Calling NoteOn with values: %+v", call.Args.NoteOn)
		var err error
		if !s.DeviceOutOpened {
			err = errors.New("can't send command when device not opened")
		} else {
			err = writer.NoteOn(s.Writer, call.Args.NoteOn.Pitch, call.Args.NoteOn.Velocity)
		}
		if err != nil {
			return err
		}
		return nil
	case types.MidiNoteOffType():
		var err error
		if !s.DeviceOutOpened {
			err = errors.New("can't send command when device not opened")
		} else {
			err = writer.NoteOff(s.Writer, call.Args.NoteOff)
		}
		if err != nil {
			return err
		}
		return nil
	default:
		return nil
	}
}
