package main

import (
	"context"

	log "github.com/sirupsen/logrus"

	"github.com/geomyidia/erl-midi-server/internal/app"
	"github.com/geomyidia/erl-midi-server/internal/cli"
	"github.com/geomyidia/erl-midi-server/pkg/erl"
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
	key := types.ParserKey("key")
	ctx := context.WithValue(context.Background(), key, flags.Parser)
	cmd := ""
	if len(flags.Args) > 0 {
		cmd = flags.Args[0]
	}
	if flags.Daemon || flags.Parser != cli.TextParser {
		server.Serve(ctx, key, flags.Parser)
	} else {
		log.Debug("using CLI mode ...")
		server.ProcessCommand(ctx, key, erl.Result(cmd))
	}
}
