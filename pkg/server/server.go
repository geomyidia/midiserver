package server

import (
	"context"
	"sync"
	"syscall"

	log "github.com/sirupsen/logrus"

	"github.com/ut-proj/midiserver/internal/util"
	"github.com/ut-proj/midiserver/pkg/commands"
	"github.com/ut-proj/midiserver/pkg/erl"
	"github.com/ut-proj/midiserver/pkg/erl/messages"
	"github.com/ut-proj/midiserver/pkg/midi"
	"github.com/ut-proj/midiserver/pkg/types"
)

func Serve(ctx context.Context, midiSys *midi.System, flags *types.Flags) {
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
		ProcessMessages(ctx, midiSys, opts, flags)
	}()

	// Listen for the interrupt signal.
	<-ctx.Done()
	// Restore default behavior on the interrupt signal and notify user of shutdown.
	cancel()
	log.Info("shutting down gracefully, press Ctrl+C again to force")
	log.Info("waiting for wait groups to finish ...")
	midiSys.Shutdown()
	wg.Wait()
	log.Info("application shutdown complete.")
}

func ProcessMessages(ctx context.Context, midiSys *midi.System, opts *erl.Opts, flags *types.Flags) {
	log.Info("processing messages sent to Go language server ...")
	log.Debugf("using command processor options %#v", opts)
	go func() {
		for {
			ProcessMessage(ctx, midiSys, opts, flags)
			continue
		}
	}()
	<-ctx.Done()
}

func ProcessMessage(ctx context.Context, midiSys *midi.System, opts *erl.Opts, flags *types.Flags) {
	mp, err := messages.NewMessageProcessor(opts)
	if err != nil {
		log.Error(err)
		return
	}
	result := mp.Process()
	if result == erl.Continue() {
		return
	}
	log.Trace("Got message type: ", result)
	if mp.IsMidi {
		callGroup := mp.MidiCallGroup()
		midiSys.Dispatch(ctx, callGroup.Calls(), callGroup.IsParallel(), flags)
	} else {
		commands.Dispatch(ctx, result.ToCommand(), mp.CommandArgs(), flags)
	}
	log.Trace("message processing complete")
}
