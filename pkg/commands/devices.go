package commands

import (
	"fmt"

	"github.com/geomyidia/midiserver/pkg/midi"
)

func ListDevices() {
	midiSystem := midi.NewSystem()
	defer midiSystem.Close()

	fmt.Printf("MIDI IN Ports:\n")
	for _, port := range midiSystem.DevicesIn {
		fmt.Printf("\t[%v] %s\n", port.Number(), port.String())
	}

	fmt.Printf("MIDI OUT Ports:\n")
	for _, port := range midiSystem.DevicesOut {
		fmt.Printf("\t[%v] %s\n", port.Number(), port.String())
	}
}
