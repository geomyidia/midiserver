package server

import (
	"context"
	"sync"
	"syscall"

	"github.com/ergo-services/ergo/gen"
	log "github.com/sirupsen/logrus"

	"github.com/ut-proj/midiserver/internal/util"
	"github.com/ut-proj/midiserver/pkg/commands"
	"github.com/ut-proj/midiserver/pkg/erl"
	"github.com/ut-proj/midiserver/pkg/erl/messages"
	"github.com/ut-proj/midiserver/pkg/erl/packets"
	"github.com/ut-proj/midiserver/pkg/erl/rpc"
	"github.com/ut-proj/midiserver/pkg/midi"
	"github.com/ut-proj/midiserver/pkg/types"
)

type GenServer struct {
	gen.Server
}

func Serve(ctx context.Context, midiSys *midi.System, flags *types.Flags) {
	log.Info("starting the server ...")
	ctx, cancel := util.SignalWithContext(ctx, syscall.SIGINT, syscall.SIGTERM, syscall.SIGKILL)
	defer cancel()
	var wg sync.WaitGroup

	wg.Add(1)
	go func() {
		defer wg.Done()
		opts := &erl.Opts{IsHexEncoded: true}
		if flags.Parser == types.PortParser() {
			opts = erl.DefaultOpts()
		}
		HandleMessages(ctx, midiSys, opts, flags)
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		rpcClient, err := rpc.New(flags)
		if err != nil {
			log.Error(err)
			cancel()
		}
		err = ReceiveMIDI(ctx, midiSys, rpcClient, flags)
		if err != nil {
			log.Error(err)
			cancel()
		}
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

func HandleMessages(ctx context.Context, midiSys *midi.System, opts *erl.Opts, flags *types.Flags) {
	log.Info("processing messages sent to Go language server ...")
	log.Debugf("using command processor options %#v", opts)
	go func() {
		for {
			HandleMessage(ctx, midiSys, opts, flags)
			continue
		}
	}()
	<-ctx.Done()
}

func HandleMessage(ctx context.Context, midiSys *midi.System, opts *erl.Opts, flags *types.Flags) {
	var resp *messages.Response
	term, err := packets.ToTerm(opts)
	if err != nil {
		log.Error(err)
		resp, _ = messages.NewResponse(types.EmptyResult, types.Err(err.Error()))
		resp.Send()
		return
	}
	log.Tracef("got Erlang ports term: %#v", term)
	msg, err := messages.New(term)
	if err != nil {
		log.Error(err)
		resp, _ = messages.NewResponse(types.EmptyResult, types.Err(err.Error()))
		resp.Send()
		return
	}

	msgName := msg.Name()
	log.Tracef("Got message name %s", msgName)
	switch msg.Type() {
	case string(types.MidiKey):
		err := midi.HandleMessage(msg.Args)
		if err != nil {
			log.Error(err)
			resp, _ = messages.NewResponse(types.EmptyResult, types.Err(err.Error()))
			resp.Send()
			return
		}
		// callGroup := mp.MidiCallGroup()
		// midiSys.Dispatch(ctx, callGroup.Calls(), callGroup.IsParallel(), flags)
		log.Debug("TODO: update MIDI message handling")
	case string(types.CommandKey):
		commands.Dispatch(ctx, msg, flags)
	default:
		err = ErrUnsupMessageType
		log.Error(err)
	}
	log.Trace("message handling complete")
}

func ReceiveMIDI(ctx context.Context, midiSys *midi.System, rpcClient *rpc.Client, flags *types.Flags) error {
	if flags.MidiInDeviceID < 0 {
		log.Warn("no valid device ID set for mini-in; not starting listener ...")
	}
	midiSys.SetReader(rpcClient, uint8(flags.MidiInDeviceID))
	log.Infof("awaiting MIDI messages on device %v ...", flags.MidiInDeviceID)
	return midiSys.Reader.ListenTo(midiSys.DeviceIn)
}
