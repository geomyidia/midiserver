package commands

import (
	"time"

	log "github.com/sirupsen/logrus"
	"gitlab.com/gomidi/midi/writer"
	driver "gitlab.com/gomidi/rtmididrv"

	"github.com/ut-proj/midiserver/pkg/types"
)

type Opts struct {
	DeviceId    uint8
	MidiChannel uint8
	Pitch       uint8
	Velocity    uint8
	Duration    uint8
}

func DefaultOpts() *Opts {
	return &Opts{
		DeviceId:    0,
		MidiChannel: 0,
		Pitch:       24,
		Velocity:    100,
		Duration:    4,
	}
}

func PlayNote(arg types.PropList) {
	var opts *Opts
	if arg == nil || len(arg) == 0 {
		log.Debug("got nil args ...")
		opts = DefaultOpts()
	} else {
		opts = &Opts{
			DeviceId:    arg["device"].(uint8),
			MidiChannel: arg["channel"].(uint8),
			Pitch:       arg["pitch"].(uint8),
			Velocity:    arg["velocity"].(uint8),
			Duration:    arg["duration"].(uint8),
		}
	}
	log.Debugf("Got opts: %+v", opts)
	drv, err := driver.New()
	if err != nil {
		log.Panic(err)
	}
	defer drv.Close()

	ins, err := drv.Ins()
	if err != nil {
		log.Panic(err)
	}

	outs, err := drv.Outs()
	if err != nil {
		log.Panic(err)
	}

	in, out := ins[0], outs[opts.DeviceId]

	err = in.Open()
	if err != nil {
		log.Panic(err)
	}

	err = out.Open()
	if err != nil {
		log.Panic(err)
	}

	wr := writer.New(out)
	wr.SetChannel(opts.MidiChannel)

	seconds := opts.Duration
	log.Debugf("playing note for %d seconds ...", seconds)

	err = writer.NoteOn(wr, opts.Pitch, opts.Velocity)
	if err != nil {
		log.Panic(err)
	}

	duration := time.Duration(seconds) * time.Second
	time.Sleep(duration)
	writer.NoteOff(wr, opts.Pitch)
}
