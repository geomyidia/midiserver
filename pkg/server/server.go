package server

import (
	"context"
	"sync"
	"syscall"

	"github.com/geomyidia/erl-midi-server/internal/util"
	"github.com/geomyidia/erl-midi-server/pkg/erl"
	"github.com/geomyidia/erl-midi-server/pkg/types"
	log "github.com/sirupsen/logrus"
)

func Serve(ctx context.Context, flags *types.Flags) {
	log.Info("starting the server ...")
	ctx, cancel := util.SignalWithContext(ctx, syscall.SIGINT, syscall.SIGTERM)
	defer cancel()
	var wg sync.WaitGroup

	wg.Add(1)
	go func() {
		defer wg.Done()
		opts := &erl.Opts{IsHexEncoded: true}
		if flags.Parser == types.PortParser() {
			opts = erl.DefaultOpts()
		}
		ProcessMessages(ctx, ProcessCommand, opts, flags)
	}()

	// Listen for the interrupt signal.
	<-ctx.Done()
	// Restore default behavior on the interrupt signal and notify user of shutdown.
	cancel()
	log.Info("shutting down gracefully, press Ctrl+C again to force")
	log.Info("waiting for wait groups to finish ...")
	wg.Wait()
	log.Info("application shutdown complete.")
}
