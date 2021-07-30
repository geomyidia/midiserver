package app

import (
	"context"
	"sync"
	"syscall"

	"github.com/geomyidia/erl-midi-server/pkg/port"
	"github.com/geomyidia/erl-midi-server/pkg/server"
	"github.com/geomyidia/erl-midi-server/pkg/types"
	log "github.com/sirupsen/logrus"
)

func Serve(ctx context.Context, key types.ParserKey, parserFlag string) {
	log.Info("Starting the server ...")
	ctx, cancel := SignalWithContext(ctx, syscall.SIGINT, syscall.SIGTERM)
	defer cancel()
	var wg sync.WaitGroup

	wg.Add(1)
	go func() {
		defer wg.Done()
		parser := port.ProcessExecMessage
		if parserFlag == PortParser {
			parser = port.ProcessPortMessage
		}
		port.ProcessMessages(ctx, parser, server.ProcessCommand, key)
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
