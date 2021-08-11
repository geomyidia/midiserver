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

	"github.com/geomyidia/midiserver/pkg/types"
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

func (s *System) SetChannel(channelId uint8) error {
	log.Trace("setting channel ...")
	s.Writer.SetChannel(channelId)
	s.ChannelSet = true
	return nil
}

func (s *System) GetChannel() channel.Channel {
	return channel.Channel(s.Writer.Channel())
}

func (s *System) Dispatch(ctx context.Context, calls []types.MidiCall, flags *types.Flags) {
	log.Debug("dispatching MIDI operation ...")
	// log.Tracef("got MIDI calls: %v", calls)
	for _, call := range calls {
		log.Debugf("making MIDI call %+v ...", call)
		err := s.CallMidi(call)
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
	if !s.ChannelSet {
		return fmt.Errorf("message of type %s require the channel to be set",
			call.Op)
	}
	ch := s.GetChannel()
	switch call.Op {
	case types.MidiNoteOnType():
		// XXX go back to the writer.NoteOn usage for this ...
		log.Tracef("calling NoteOn with values: %+v", call.Args.NoteOn)
		var err error
		if !s.DeviceOutOpened {
			return errors.New("can't send command when device not opened")
		} else {
			msg := ch.NoteOn(call.Args.NoteOn.Pitch, call.Args.NoteOn.Velocity)
			log.Tracef("created MIDI msg: %+v", msg)
			err = s.Writer.Write(msg)
			if err != nil {
				return err
			}
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
		log.Errorf("no handler for operation '%s'", call.Op)
		return nil
	}
}
