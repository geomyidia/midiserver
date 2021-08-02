package main

import (
	"context"

	log "github.com/sirupsen/logrus"

	"github.com/geomyidia/erl-midi-server/internal/app"
	"github.com/geomyidia/erl-midi-server/internal/cli"
	"github.com/geomyidia/erl-midi-server/pkg/server"
	"github.com/geomyidia/erl-midi-server/pkg/types"
	"github.com/geomyidia/erl-midi-server/pkg/version"
)

func main() {
	flags := cli.Parse()
	app.SetupLogging(flags.LogLevel)
	log.Info("Welcome to the Go midiserver!")
	log.Infof("running version: %s", version.VersionedBuildString())
	log.Tracef("flags: %+v", flags)
	app.SetupRandom()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	if flags.Daemon || (flags.Parser != types.TextParser()) {
		server.Serve(ctx, flags)
	} else {
		log.Debug("using CLI mode ...")
		server.ProcessCommand(ctx, flags.Command, flags)
	}
}
