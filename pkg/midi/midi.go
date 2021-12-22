package midi

import (
	"context"
	"errors"
	"fmt"

	log "github.com/sirupsen/logrus"
	"gitlab.com/gomidi/midi"
	"gitlab.com/gomidi/midi/midimessage/channel"
	"gitlab.com/gomidi/midi/reader"
	"gitlab.com/gomidi/midi/writer"
	"gitlab.com/gomidi/rtmididrv"

	"github.com/ut-proj/midiserver/pkg/types"
)

const (
	MSBBankCC uint8 = 0
	LSBBankCC uint8 = 32
)

type System struct {
	Driver          midi.Driver
	DevicesIn       []midi.In
	DevicesOut      []midi.Out
	DeviceIn        midi.In
	DeviceOut       midi.Out
	Reader          *reader.Reader
	Writer          *writer.Writer
	DeviceInOpened  bool
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
		DeviceInOpened:  false,
		DevicesOut:      outs,
		DeviceOutOpened: false,
		ChannelSet:      false,
	}
}

func (s *System) Shutdown() {
	log.Info("shutting down MIDI system ...")
	s.Driver.Close()
	if s.DeviceOut != nil && s.DeviceOut.IsOpen() {
		s.DeviceOut.Close()
	}
	s.DeviceOutOpened = false
	if s.DeviceIn != nil && s.DeviceIn.IsOpen() {
		s.DeviceIn.Close()
	}
	s.DeviceInOpened = false
}

func (s *System) SetWriter(deviceOutID uint8) error {
	if s.DeviceOutOpened {
		return nil
	}
	log.Info("setting device ...")
	s.DeviceOut = s.DevicesOut[deviceOutID]
	err := s.DeviceOut.Open()
	if err != nil {
		log.Fatal(err)
		return err
	}
	s.Writer = writer.New(s.DeviceOut)
	s.DeviceOutOpened = true
	return nil
}

func (s *System) SetReader(deviceInID uint8) {
	s.Reader = reader.New(
		reader.NoLogger(),
		reader.Each(ReceiveEach),
		reader.RTClock(ReceiveClock),
		reader.RTContinue(ReceiveContinue),
		reader.RTReset(ReceiveReset),
		reader.RTStart(ReceiveStart),
		reader.RTStop(ReceiveStop),
		reader.RTTick(ReceiveTick),
		reader.Unknown(ReceiveUnknown),
	)
	s.DeviceIn = s.DevicesIn[deviceInID]
	s.DeviceIn.Open()
	s.DeviceInOpened = true
}

func (s *System) SetWriterChannel(channelId uint8) error {
	log.Info("setting channel ...")
	s.Writer.SetChannel(channelId)
	s.ChannelSet = true
	log.Tracef("current channel value: %v", s.GetChannel())
	return nil
}

func (s *System) GetWriterChannel() channel.Channel {
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
		return s.SetWriter(call.Args.Device)
	case types.MidiChannelType():
		return s.SetWriterChannel(call.Args.Channel)
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
	case types.MidiProgramChangeType():
		err := writer.ProgramChange(s.Writer, call.Args.Program)
		if err != nil {
			return err
		}
		return nil
	case types.MidiBankSelectMSBType():
		err := writer.ControlChange(s.Writer, MSBBankCC, call.Args.CC.Value)
		if err != nil {
			return err
		}
		return nil
	case types.MidiBankSelectLSBType():
		err := writer.ControlChange(s.Writer, LSBBankCC, call.Args.CC.Value)
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
	case types.MidiRealtimeType():
		log.Tracef("got RT message with args: %+v", call.Args)
		switch call.Args.Realtime {
		case types.MidiRTClock():
			err := writer.RTClock(s.Writer)
			if err != nil {
				return err
			}
			return nil
		case types.MidiRTContinue():
			err := writer.RTContinue(s.Writer)
			if err != nil {
				return err
			}
			return nil
		case types.MidiRTReset():
			err := writer.RTReset(s.Writer)
			if err != nil {
				return err
			}
			return nil
		case types.MidiRTStart():
			err := writer.RTStart(s.Writer)
			if err != nil {
				return err
			}
			return nil
		case types.MidiRTStop():
			err := writer.RTStop(s.Writer)
			if err != nil {
				return err
			}
			return nil
		case types.MidiRTTick():
			err := writer.RTTick(s.Writer)
			if err != nil {
				return err
			}
			return nil
		default:
			log.Errorf("unsupported realtime message '%s'", call.Args.Realtime)
			return nil
		}
	default:
		log.Errorf("no handler for operation '%s'", call.Op)
		return nil
	}
}

// DEPRECATED FUNCTIONS

func (s *System) SetDevice(deviceId uint8) error {
	log.Warnf("the 'SetDevice' method is deprecated; use 'SetWriter' instead")
	return s.SetWriter(deviceId)
}

func (s *System) SetChannel(channelId uint8) error {
	log.Warnf("the 'SetChannel' method is deprecated; use 'SetWriterChannel' instead")
	return s.SetWriterChannel(channelId)
}

func (s *System) GetChannel() channel.Channel {
	log.Warnf("the 'GetChannel' method is deprecated; use 'GetWriterChannel' instead")
	return s.GetWriterChannel()
}
