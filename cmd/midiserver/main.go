package main

import (
	"context"

	log "github.com/sirupsen/logrus"

	"github.com/geomyidia/midiserver/internal/app"
	"github.com/geomyidia/midiserver/internal/cli"
	"github.com/geomyidia/midiserver/pkg/commands"
	"github.com/geomyidia/midiserver/pkg/server"
	"github.com/geomyidia/midiserver/pkg/types"
	"github.com/geomyidia/midiserver/pkg/version"
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
		// XXX fill this up
		args := make(types.Proplist)
		commands.Dispatch(ctx, flags.Command, args, flags)
	}
}
