package main

import (
	"context"

	log "github.com/sirupsen/logrus"

	"github.com/ut-proj/midiserver/internal/app"
	"github.com/ut-proj/midiserver/internal/cli"
	"github.com/ut-proj/midiserver/pkg/commands"
	"github.com/ut-proj/midiserver/pkg/midi"
	"github.com/ut-proj/midiserver/pkg/server"
	"github.com/ut-proj/midiserver/pkg/types"
)

func main() {
	flags := cli.Parse()
	app.Setup(flags)
	midiSystem := midi.NewSystem()
	defer midiSystem.Shutdown()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	if flags.Daemon || (flags.Parser != types.TextParser()) {
		server.Serve(ctx, midiSystem, flags)
	} else {
		log.Debug("using CLI mode ...")
		args := make(map[string]interface{})
		commands.Dispatch(ctx, flags.Command, args, flags)
	}
}
