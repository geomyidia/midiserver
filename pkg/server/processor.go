package server

import (
	"context"

	log "github.com/sirupsen/logrus"

	"github.com/geomyidia/erl-midi-server/pkg/erl"
	"github.com/geomyidia/erl-midi-server/pkg/types"
)

// CommandProcessor ...
type CommandProcessor func(context.Context, types.ParserKey, erl.Result)
type MessageProcessor func() erl.Result

func ProcessMessage(opts *erl.Opts) erl.Result {
	log.Debug("process message ...")
	mp, err := erl.NewMessageProcessor(opts)
	if err != nil {
		log.Error(err)
		return mp.Continue()
	}
	return mp.Process()
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
func ProcessMessages(ctx context.Context, cmdFn CommandProcessor, key types.ParserKey, opts *erl.Opts) {
	log.Info("processing messages sent to Go language server ...")
	log.Debugf("using command processor %T", cmdFn)
	log.Debugf("using command processor options %#v", opts)
	go func() {
		for {
			cmd := ProcessMessage(opts)
			if cmd == erl.Continue() {
				continue
			}
			cmdFn(ctx, key, cmd)

		}
	}()
	<-ctx.Done()
}
