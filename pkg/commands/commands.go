package commands

import (
	"context"
	"fmt"
	"strconv"

	log "github.com/sirupsen/logrus"

	"github.com/ut-proj/midiserver/pkg/erl/messages"
	"github.com/ut-proj/midiserver/pkg/text"
	"github.com/ut-proj/midiserver/pkg/types"
	"github.com/ut-proj/midiserver/pkg/version"
)

// Dispatch ...
func Dispatch(
	ctx context.Context,
	command types.CommandType,
	args map[string]interface{},
	flags *types.Flags,
) {
	log.Debug("Dispatching command ...")
	var result types.Result
	var err types.Err
	switch command {
	case types.PingCommand():
		if flags.RemoteNode != "" && flags.RemoteModule != "" {
			err := PingRemoteModule(flags)
			if err != nil {
				log.Error(err)
			} else {
				result = types.Result("ok")
			}
		} else {
			result = types.Result("pong")
		}
	case types.PlayNoteCommand():
		if len(flags.Args) == 6 {
			args["device"] = toUint(flags.Args[1])
			args["channel"] = toUint(flags.Args[2])
			args["pitch"] = toUint(flags.Args[3])
			args["velocity"] = toUint(flags.Args[4])
			args["duration"] = toUint(flags.Args[5])
		}
		PlayNote(args)
		result = types.Result("ok")
	case types.ExampleCommand():
		if len(flags.Args) == 3 {
			args["device"] = toUint(flags.Args[1])
			args["channel"] = toUint(flags.Args[2])
		}
		PlayExample(args)
		result = types.Result("ok")
	case types.ListDevicesCommand():
		ListDevices()
		result = types.Result("ok")
	case types.ListNodesCommand():
		ListNodes(flags)
		result = types.Result("ok")
	case types.RemotePortCommand():
		ShowRemotePort(flags)
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
		resp, _ := messages.NewResponse(result, err)
		resp.Send()
	} else if flags.Parser == types.TextParser() {
		resp := text.NewResponse(result, err)
		resp.Send()
	} else {
		log.Errorf("unsupported parser type: %v", flags.Parser)
	}
}

func toUint(v string) uint8 {
	r, err := strconv.ParseUint(v, 10, 8)
	if err != nil {
		log.Error(err)
	}
	return uint8(r)
}
