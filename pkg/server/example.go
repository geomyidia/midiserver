package server

import (
	"time"

	log "github.com/sirupsen/logrus"
	"gitlab.com/gomidi/midi/writer"
	driver "gitlab.com/gomidi/rtmididrv"
)

func Example() {

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

	in, out := ins[0], outs[0]

	err = in.Open()
	if err != nil {
		log.Panic(err)
	}

	err = out.Open()
	if err != nil {
		log.Panic(err)
	}

	wr := writer.New(out)
	wr.SetChannel(0) // sets the channel for the next messages

	seconds := 5
	log.Debugf("Playing note for %d seconds ...", seconds)

	err = writer.NoteOn(wr, 60, 100)
	if err != nil {
		log.Panic(err)
	}

	duration := time.Duration(seconds) * time.Second
	time.Sleep(duration)
	writer.NoteOff(wr, 60) // let the note ring for 5 sec
}
