package midiserver

import (
	"context"

	log "github.com/sirupsen/logrus"

	"github.com/geomyidia/erl-midi-server/pkg/port"
	"github.com/geomyidia/erl-midi-server/pkg/types"
	"github.com/geomyidia/erl-midi-server/pkg/version"
)

// ProcessCommand ...
func ProcessCommand(ctx context.Context, key types.ParserKey, command string) {
	parserType := ctx.Value(key).(string)
	switch command {
	case "ping":
		sendResult(parserType, "pong")
	case "example":
		Example()
		sendResult(parserType, "ok")
	case "stop":
		log.Info("Stopping Go MIDI server ...")
		<-ctx.Done()
	case "version":
		sendResult(parserType, version.VersionedBuildString())
	default:
		sendError(parserType, "Received unsupported command: "+command)
	}
}

func sendResult(parserType string, msg string) {
	if parserType == "text" {
		println(msg)
	} else {
		port.SendResult(msg)
	}
}

func sendError(parserType string, msg string) {
	if parserType == "text" {
		log.Error(msg)
	} else {
		port.SendError(msg)
	}
}
