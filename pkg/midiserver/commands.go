package midiserver

import (
	"context"

	log "github.com/sirupsen/logrus"

	"github.com/geomyidia/erl-midi-server/pkg/port"
	"github.com/geomyidia/erl-midi-server/pkg/version"
)

// ProcessCommand ...
func ProcessCommand(ctx context.Context, command string) {
	switch command {
	case "ping":
		port.SendResult("pong")
	case "example":
		Example()
		port.SendResult("ok")
	case "stop":
		log.Info("Stopping Go MIDI server ...")
		<-ctx.Done()
	case "version":
		port.SendResult(version.VersionedBuildString())
	default:
		port.SendError("Received unsupported command: " + command)
	}
}
