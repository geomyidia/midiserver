package commands

import (
	"fmt"

	"github.com/geomyidia/midiserver/pkg/midi"
)

func ListDevices() {
	midiSystem := midi.NewSystem()
	defer midiSystem.Shutdown()

	fmt.Printf("MIDI IN devices:\n")
	for _, port := range midiSystem.DevicesIn {
		fmt.Printf("\t[%v] %s\n", port.Number(), port.String())
	}

	fmt.Printf("MIDI OUT devices:\n")
	for _, port := range midiSystem.DevicesOut {
		fmt.Printf("\t[%v] %s\n", port.Number(), port.String())
	}
}
