package midi

import (
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

func Dispatch(data interface{}, flags *types.Flags) {
	// Process the new messages defined here:
	// * https://github.com/erlsci/midilib/blob/release/0.4.x/src/midimsg.erl
	// that have been processed in pkg/erl/...?
}
