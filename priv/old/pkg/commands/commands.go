package commands

import (
	"context"
	"fmt"
	"strconv"

	"github.com/ergo-services/ergo/etf"
	log "github.com/sirupsen/logrus"

	"github.com/geomyidia/erlcmd/pkg/messages"

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
	var result messages.Result
	var err messages.Err
	if msg.Type() != string(messages.CommandKey) {
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
				result = messages.OkResult
			}
		} else {
			result = messages.PongResult
		}
	case types.PlayNoteCommand():
		// TODO: let's put this logic into something that parses the
		//       playnote message and creates a dedicated struct for
		//       it
		args := make(map[string]interface{})
		for _, arg := range msg.Args() {
			tuple, ok := arg.(etf.Tuple)
			if !ok {
				log.Error(messages.ErrMsgTupleFormat)
				continue
			}
			args[tuple.Element(1).(string)] = tuple.Element(2).(uint)
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
		result = messages.OkResult
	case types.ExampleCommand():
		args := make(map[string]interface{})
		for _, arg := range msg.Args() {
			tuple, ok := arg.(etf.Tuple)
			if !ok {
				log.Error(messages.ErrMsgTupleFormat)
				continue
			}
			args[tuple.Element(1).(string)] = tuple.Element(2).(uint)
		}
		// TODO: if the values aren't in the payload, pull them
		// from the flags
		// if len(flags.Args) == 3 {
		// 	args["device"] = toUint(flags.Args[1])
		// 	args["channel"] = toUint(flags.Args[2])
		// }
		PlayExample(args)
		result = messages.OkResult
	case types.ListDevicesCommand():
		ListDevices()
		result = messages.OkResult
	case types.ListNodesCommand():
		ListNodes(flags)
		result = messages.OkResult
	case types.RemotePortCommand():
		ShowRemotePort(flags)
		result = messages.OkResult
	case types.StopCommand():
		log.Info("stopping Go MIDI server ...")
		result = messages.StoppingResult
		<-ctx.Done()
	case types.VersionCommand():
		result = messages.Result(version.VersionedBuildString())
	case types.EmptyCommand():
		result = messages.Result("missing command; see -h for useage")
	default:
		result = messages.Result(
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
