package server

import (
	"context"

	log "github.com/sirupsen/logrus"

	"github.com/geomyidia/midiserver/pkg/commands"
	"github.com/geomyidia/midiserver/pkg/erl"
	"github.com/geomyidia/midiserver/pkg/types"
)

func ProcessMessage(ctx context.Context, opts *erl.Opts, flags *types.Flags) {
	mp, err := erl.NewMessageProcessor(opts)
	if err != nil {
		log.Error(err)
		return
	}
	result := mp.Process()
	if result == erl.Continue() {
		return
	}
	log.Warning(result)
	log.Debugf("Got MIDI data: %+v", mp.Midi())
	commands.Dispatch(ctx, result.ToCommand(), mp.CommandArgs(), flags)
	log.Debug("processed message ...")
	return
}

// ProcessMessages handles messages of the Erlang Port format along the
// following lines:
//   a           = []byte{0x83, 0x64, 0x0, 0x1, 0x61, 0xa}
//   "a"         = []byte{0x83, 0x6b, 0x0, 0x1, 0x61, 0xa}
//   {}          = []byte{0x83, 0x68, 0x0, 0xa}
//   {a}         = []byte{0x83, 0x68, 0x1, 0x64, 0x0, 0x1, 0x61, 0xa}
//   {"a"}       = []byte{0x83, 0x68, 0x1, 0x6b, 0x0, 0x1, 0x61, 0xa}
//   {a, a}      = []byte{0x83, 0x68, 0x2, 0x64, 0x0, 0x1, 0x61, 0x64, 0x0, 0x1, 0x61, 0xa}
//   {a, test}   = []byte{0x83, 0x68, 0x2, 0x64, 0x0, 0x1, 0x61, 0x64, 0x0, 0x4, 0x74, 0x65, 0x73, 0x74, 0xa}
//   {a, "test"} = []byte{0x83, 0x68, 0x2, 0x64, 0x0, 0x1, 0x61, 0x6b, 0x0, 0x4, 0x74, 0x65, 0x73, 0x74, 0xa}
func ProcessMessages(ctx context.Context, opts *erl.Opts, flags *types.Flags) {
	log.Info("processing messages sent to Go language server ...")
	log.Debugf("using command processor options %#v", opts)
	go func() {
		for {
			ProcessMessage(ctx, opts, flags)
			continue
		}
	}()
	<-ctx.Done()
}
