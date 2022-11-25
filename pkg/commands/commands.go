package commands

import (
	"context"
	"fmt"
	"strconv"

	"github.com/ut-proj/midiserver/pkg/erl/datatypes"

	log "github.com/sirupsen/logrus"

	"github.com/ut-proj/midiserver/pkg/erl/messages"
	"github.com/ut-proj/midiserver/pkg/text"
	"github.com/ut-proj/midiserver/pkg/types"
	"github.com/ut-proj/midiserver/pkg/version"
)

// Dispatch ...
func Dispatch(
	ctx context.Context,
	msg *messages.Message,
	flags *types.Flags,
) {
	log.Debug("Dispatching command ...")
	var result types.Result
	var err types.Err
	if msg.Type() != string(types.CommandKey) {
		log.Error(ErrCmdMsgFormat)
		resp := text.NewResponse(result, err)
		resp.Send()
	}
	cmdName := msg.Name()
	switch types.Command(types.CommandName(cmdName)) {
	case types.PingCommand():
		if flags.RemoteNode != "" && flags.RemoteModule != "" {
			err := PingRemoteModule(flags)
			if err != nil {
				log.Error(err)
			} else {
				result = types.OkResult
			}
		} else {
			result = types.PongResult
		}
	case types.PlayNoteCommand():
		// TODO: let's put this logic into something that parses the
		//       playnote message and creates a dedicated struct for
		//       it
		args := make(map[string]interface{})
		for _, arg := range msg.Args() {
			tuple, ok := arg.(*datatypes.Tuple)
			if !ok {
				log.Error(datatypes.ErrCastingTuple)
				continue
			}
			args[tuple.Key().(string)] = tuple.Value().(uint)
		}
		// TODO: if the values aren't in the payload, pull them
		// from the flags
		// deviceIdx = 1
		// channelIdx = 2
		// pitchIdx = 3
		// velocityIdx = 4
		// durationIdx = 5
		// if len(flags.Args) == 6 {
		// 	args["device"] = toUint(flags.Args[deviceIdx])
		// 	args["channel"] = toUint(flags.Args[channelIdx])
		// 	args["pitch"] = toUint(flags.Args[pitchIdx])
		// 	args["velocity"] = toUint(flags.Args[velocityIdx])
		// 	args["duration"] = toUint(flags.Args[durationIdx])
		// }
		PlayNote(args)
		result = types.OkResult
	case types.ExampleCommand():
		args := make(map[string]interface{})
		for _, arg := range msg.Args() {
			tuple, ok := arg.(*datatypes.Tuple)
			if !ok {
				log.Error(datatypes.ErrCastingTuple)
				continue
			}
			args[tuple.Key().(string)] = tuple.Value().(uint)
		}
		// TODO: if the values aren't in the payload, pull them
		// from the flags
		// if len(flags.Args) == 3 {
		// 	args["device"] = toUint(flags.Args[1])
		// 	args["channel"] = toUint(flags.Args[2])
		// }
		PlayExample(args)
		result = types.OkResult
	case types.ListDevicesCommand():
		ListDevices()
		result = types.OkResult
	case types.ListNodesCommand():
		ListNodes(flags)
		result = types.OkResult
	case types.RemotePortCommand():
		ShowRemotePort(flags)
		result = types.OkResult
	case types.StopCommand():
		log.Info("stopping Go MIDI server ...")
		result = types.StoppingResult
		<-ctx.Done()
	case types.VersionCommand():
		result = types.Result(version.VersionedBuildString())
	case types.EmptyCommand():
		result = types.Result("missing command; see -h for useage")
	default:
		result = types.Result(
			fmt.Sprintf(
				"received unsupported command: '%s' (data: %+v, type: %T)",
				cmdName, msg, msg))
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
