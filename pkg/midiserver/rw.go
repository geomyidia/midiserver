package midiserver

import (
	"fmt"
	"io"
	"time"

	log "github.com/sirupsen/logrus"
	driver "gitlab.com/gomidi/rtmididrv"
	"gitlab.com/gomidi/midi/reader"
	"gitlab.com/gomidi/midi/writer"
)

type printer struct{}

func (pr printer) noteOn(p *reader.Position, channel, pitch, vel uint8) {
	fmt.Printf("NoteOn (ch %v: pitch %v vel: %v)\n", channel, pitch, vel)
}

func (pr printer) noteOff(p *reader.Position, channel, pitch, vel uint8) {
	fmt.Printf("NoteOff (ch %v: pitch %v)\n", channel, pitch)
}

func Example() {

	var p printer

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

	log.Debug("MIDI IN Ports")
	for _, port := range ins {
		log.Debugf("[%v] %s\n", port.Number(), port.String())
	}

	log.Debug("MIDI OUT Ports")
	for _, port := range outs {
		log.Debugf("[%v] %s\n", port.Number(), port.String())
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

	// to disable logging, pass mid.NoLogger() as option
	rd := reader.New(reader.NoLogger(),
		// set the callbacks for the messages you are interested in
		reader.NoteOn(p.noteOn),
		reader.NoteOff(p.noteOff),
	)

	// to allow reading and writing concurrently in this example
	// we need a pipe
	piperd, pipewr := io.Pipe()

	go func() {
		// wr := writer.New(pipewr)
		wr := writer.New(out)
		wr.SetChannel(0) // sets the channel for the next messages
		err := writer.NoteOn(wr, 60, 100)
		if err != nil {
			log.Panic(err)
		}
		time.Sleep(5 * time.Second)
		writer.NoteOff(wr, 60) // let the note ring for 5 sec
		pipewr.Close()          // finishes the writing
	}()

	for {
		if reader.ReadAllFrom(rd, piperd) == io.EOF {
			piperd.Close() // finishes the reading
			break
		}
	}
}
