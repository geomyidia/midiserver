package main

import (
	log "github.com/sirupsen/logrus"

	"github.com/geomyidia/erl-midi-server/internal/app"
	"github.com/geomyidia/erl-midi-server/pkg/version"
)

func main() {
	flags := app.ParseCLI()
	app.SetupLogging(flags.LogLevel)
	log.Info("Welcome to the Go midiserver!")
	log.Infof("Running version: %s", version.VersionedBuildString())
	app.SetupRandom()
	log.Tracef("Flags: %+v", flags)
	if flags.Daemon || flags.Parser != app.TextParser {
		app.Serve(flags.Parser)
	} else {
		log.Debug("Using CLI mode ...")
	}
}
