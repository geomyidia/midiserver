package commands

import (
	"context"
	"fmt"

	log "github.com/sirupsen/logrus"

	"github.com/geomyidia/midiserver/pkg/erl/messages"
	"github.com/geomyidia/midiserver/pkg/text"
	"github.com/geomyidia/midiserver/pkg/types"
	"github.com/geomyidia/midiserver/pkg/version"
)

// Dispatch ...
func Dispatch(ctx context.Context, command types.CommandType,
	args types.PropList, flags *types.Flags) {
	log.Debug("Dispatching command ...")
	var result types.Result
	var err types.Err
	switch command {
	case types.PingCommand():
		result = types.Result("pong")
	case types.ExampleCommand():
		Example(args)
		result = types.Result("ok")
	case types.ListDevicesCommand():
		ListDevices()
		result = types.Result("ok")
	case types.StopCommand():
		log.Info("stopping Go MIDI server ...")
		result = types.Result("stopping")
		<-ctx.Done()
	case types.VersionCommand():
		result = types.Result(version.VersionedBuildString())
	case types.EmptyCommand():
		result = types.Result("missing command; see -h for useage")
	default:
		result = types.Result(
			fmt.Sprintf("received unsupported command: '%v' (type %T)",
				command, command))
	}

	if flags.Parser == types.ExecParser() || flags.Parser == types.PortParser() {
		resp := messages.NewResponse(result, err)
		resp.Send()
	} else if flags.Parser == types.TextParser() {
		resp := text.NewResponse(result, err)
		resp.Send()
	} else {
		log.Errorf("unsupported parser type: %v", flags.Parser)
	}
}
