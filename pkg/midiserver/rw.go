package midiserver

import (
	"fmt"
	"io"
	"time"

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
		wr := writer.New(pipewr)
		wr.SetChannel(1) // sets the channel for the next messages
		writer.NoteOn(wr, 120, 100)
		time.Sleep(5 * time.Second)
		writer.NoteOff(wr, 120) // let the note ring for 5 sec
		pipewr.Close()          // finishes the writing
	}()

	for {
		if reader.ReadAllFrom(rd, piperd) == io.EOF {
			piperd.Close() // finishes the reading
			break
		}
	}
}
