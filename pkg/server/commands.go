package server

import (
	"context"
	"fmt"

	log "github.com/sirupsen/logrus"

	"github.com/geomyidia/erl-midi-server/pkg/erl/messages"
	"github.com/geomyidia/erl-midi-server/pkg/erl/term"
	"github.com/geomyidia/erl-midi-server/pkg/midi"
	"github.com/geomyidia/erl-midi-server/pkg/types"
	"github.com/geomyidia/erl-midi-server/pkg/version"
)

// ProcessCommand ...
func ProcessCommand(ctx context.Context, key types.ParserKey, command term.Result) {
	parserType := ctx.Value(key).(string)
	switch command {
	case "midi":
		midi.MessageDispatch()
	case "ping":
		sendResult(parserType, "pong")
	case "example":
		Example()
		sendResult(parserType, "ok")
	case "list-devices":
		listDevices()
	case "stop":
		log.Info("Stopping Go MIDI server ...")
		<-ctx.Done()
	case "version":
		sendResult(parserType, version.VersionedBuildString())
	default:
		if command == "" {
			command = "(no value)"
		}
		sendError(parserType, "Received unsupported command: "+string(command))
	}
}

func sendResult(parserType string, msg string) {
	if parserType == "text" {
		println(msg)
	} else {
		messages.SendResult(msg)
	}
}

func sendError(parserType string, msg string) {
	if parserType == "text" {
		log.Error(msg)
	} else {
		messages.SendError(msg)
	}
}

func listDevices() {

	midiSystem := midi.NewSystem()
	defer midiSystem.Close()

	fmt.Printf("MIDI IN Ports:\n")
	for _, port := range midiSystem.Ins {
		fmt.Printf("\t[%v] %s\n", port.Number(), port.String())
	}

	fmt.Printf("MIDI OUT Ports:\n")
	for _, port := range midiSystem.Outs {
		fmt.Printf("\t[%v] %s\n", port.Number(), port.String())
	}

}
