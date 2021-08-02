package server

import (
	"context"
	"fmt"

	log "github.com/sirupsen/logrus"

	"github.com/geomyidia/erl-midi-server/pkg/erl"
	"github.com/geomyidia/erl-midi-server/pkg/midi"
	"github.com/geomyidia/erl-midi-server/pkg/text"
	"github.com/geomyidia/erl-midi-server/pkg/types"
	"github.com/geomyidia/erl-midi-server/pkg/version"
)

// ProcessCommand ...
func ProcessCommand(ctx context.Context, command types.CommandType, flags *types.Flags) {
	var result types.Result
	var err types.Err
	switch command {
	case types.MidiCommand():
		midi.MessageDispatch()
	case types.PingCommand():
		result = types.Result("pong")
	case types.ExampleCommand():
		Example()
		result = types.Result("ok")
	case types.ListDevicesCommand():
		listDevices()
		result = types.Result("ok")
	case types.StopCommand():
		log.Info("stopping Go MIDI server ...")
		result = types.Result("stopping")
		<-ctx.Done()
	case types.VersionCommand():
		result = types.Result(version.VersionedBuildString())
	default:
		if command == "" {
			command = "(no value)"
		}
		result = types.Result("received unsupported command: " + string(command))
	}

	if flags.Parser == types.ExecParser() || flags.Parser == types.PortParser() {
		resp := erl.NewResponse(result, err)
		resp.Send()
	} else if flags.Parser == types.TextParser() {
		resp := text.NewResponse(result, err)
		resp.Send()
	} else {
		log.Errorf("unexpected parser type: %v", flags.Parser)
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
