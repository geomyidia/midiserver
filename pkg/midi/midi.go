package midi

import (
	"context"

	log "github.com/sirupsen/logrus"
	"gitlab.com/gomidi/midi"
	"gitlab.com/gomidi/rtmididrv"

	"github.com/geomyidia/midiserver/pkg/types"
)

type System struct {
	Driver *rtmididrv.Driver
	Ins    []midi.In
	Outs   []midi.Out
}

func NewSystem() *System {
	drv, err := rtmididrv.New()
	if err != nil {
		log.Error(err)
	}

	ins, err := drv.Ins()
	if err != nil {
		log.Error(err)
	}

	outs, err := drv.Outs()
	if err != nil {
		log.Error(err)
	}

	return &System{
		Driver: drv,
		Ins:    ins,
		Outs:   outs,
	}
}

func (s *System) Close() {
	s.Driver.Close()
}

func Dispatch(ctx context.Context, calls []types.MidiCall, flags *types.Flags) {
	log.Debug("Dispatching MIDI operation ...")
	log.Tracef("Got MIDI call: %+v", calls)
	// switch op {
	// case types.MidiBatch():
	// }
}
