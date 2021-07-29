package app

import (
	"context"
	"sync"
	"syscall"

	"github.com/geomyidia/erl-midi-server/pkg/midiserver"
	"github.com/geomyidia/erl-midi-server/pkg/port"
	log "github.com/sirupsen/logrus"
)

func Serve(parserFlag string) {
	log.Info("Starting the server ...")
	ctx, cancel := SignalWithContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer cancel()
	var wg sync.WaitGroup

	wg.Add(1)
	go func() {
		defer wg.Done()
		parser := port.ProcessExecMessage
		if parserFlag == PortParser {
			parser = port.ProcessPortMessage
		}
		port.ProcessMessages(ctx, parser, midiserver.ProcessCommand)
	}()

	// Listen for the interrupt signal.
	<-ctx.Done()
	// Restore default behavior on the interrupt signal and notify user of shutdown.
	cancel()
	log.Info("Shutting down gracefully, press Ctrl+C again to force")
	log.Info("Waiting for wait groups to finish ...")
	wg.Wait()
	log.Info("Application shutdown complete.")
}
