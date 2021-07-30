package main

import (
	"context"
	"flag"

	log "github.com/sirupsen/logrus"

	"github.com/geomyidia/erl-midi-server/internal/app"
	"github.com/geomyidia/erl-midi-server/pkg/server"
	"github.com/geomyidia/erl-midi-server/pkg/types"
	"github.com/geomyidia/erl-midi-server/pkg/version"
)

func main() {
	flags := app.ParseCLI()
	app.SetupLogging(flags.LogLevel)
	log.Info("Welcome to the Go midiserver!")
	log.Infof("Running version: %s", version.VersionedBuildString())
	log.Tracef("Flags: %+v", flags)
	log.Tracef("Args: %+v", flag.Args())
	app.SetupRandom()
	key := types.ParserKey("key")
	ctx := context.WithValue(context.Background(), key, flags.Parser)
	cmd := ""
	if len(flag.Args()) > 0 {
		cmd = flag.Args()[0]
	}
	if flags.Daemon || flags.Parser != app.TextParser {
		app.Serve(ctx, key, flags.Parser)
	} else {
		log.Debug("Using CLI mode ...")
		server.ProcessCommand(ctx, key, cmd)
	}
}
