package server

import (
	"time"

	log "github.com/sirupsen/logrus"
	"gitlab.com/gomidi/midi/writer"
	driver "gitlab.com/gomidi/rtmididrv"

	"github.com/geomyidia/midiserver/pkg/types"
)

type Opts struct {
	DeviceId    uint8
	MidiChannel uint8
	Pitch       uint8
	Velocity    uint8
	Duration    uint8
}

func Example(arg types.Proplist) {
	opts := &Opts{
		DeviceId:    arg["device"].(uint8),
		MidiChannel: arg["channel"].(uint8),
		Pitch:       arg["pitch"].(uint8),
		Velocity:    arg["velocity"].(uint8),
		Duration:    arg["duration"].(uint8),
	}
	log.Tracef("Got opts: %+v", opts)
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

	log.Debug("MIDI IN Ports:")
	for _, port := range ins {
		log.Debugf("\t[%v] %s", port.Number(), port.String())
	}

	log.Debug("MIDI OUT Ports:")
	for _, port := range outs {
		log.Debugf("\t[%v] %s", port.Number(), port.String())
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
	wr.SetChannel(opts.MidiChannel) // sets the channel for the next messages

	seconds := opts.Duration
	log.Debugf("playing note for %d seconds ...", seconds)

	err = writer.NoteOn(wr, opts.Pitch, opts.Velocity)
	if err != nil {
		log.Panic(err)
	}

	duration := time.Duration(seconds) * time.Second
	time.Sleep(duration)
	writer.NoteOff(wr, 60) // let the note ring for 5 sec
}
