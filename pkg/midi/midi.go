package midi

import (
	"context"

	log "github.com/sirupsen/logrus"
	"gitlab.com/gomidi/midi"
	"gitlab.com/gomidi/midi/writer"
	"gitlab.com/gomidi/rtmididrv"

	"github.com/geomyidia/midiserver/pkg/types"
)

type System struct {
	Driver     *rtmididrv.Driver
	DevicesIn  []midi.In
	DevicesOut []midi.Out
	DeviceIn   midi.In
	DeviceOut  midi.Out
	Writer     *writer.Writer
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

func (s *System) Close() {
	log.Info("Shutting down MIDI system ...")
	s.Driver.Close()
}

func (s *System) SetDevice(deviceId uint8) {
	s.DeviceOut = s.DevicesOut[deviceId]
	s.Writer = writer.New(s.DeviceOut)
}

func (s *System) SetChannel(channelId uint8) {
	s.Writer.SetChannel(channelId)
}

func (s *System) Dispatch(ctx context.Context, calls []types.MidiCall, flags *types.Flags) {
	log.Debug("Dispatching MIDI operation ...")
	log.Tracef("Got MIDI calls: %v", calls)
	for _, call := range calls {
		log.Debugf("Making MIDI call %v ...", call)
		s.CallMidi(call)
	}
}

func (s *System) CallMidi(call types.MidiCall) {
	switch call.Op {
	case types.MidiDeviceType():
		s.SetDevice(call.Args.Device)
	case types.MidiChannelType():
		s.SetChannel(call.Args.Channel)
	case types.MidiMeterType():
		println("tbd")
	case types.MidiTempoType():
		println("tbd")
	case types.MidiNoteOnType():
		println("tbd")
	case types.MidiNoteOffType():
		println("tbd")
	}
}
