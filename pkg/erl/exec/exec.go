package exec

import (
	log "github.com/sirupsen/logrus"

	"github.com/geomyidia/erl-midi-server/pkg/erl/term"
)

// ProcessMessage - when using the Erlang exec library to send
// messages to this Go server, a stange thing happens: a byte is
// dropped from the middle of the
func ProcessMessage() term.Result {
	opts := &term.Opts{IsHexEncoded: true}
	mp, err := term.NewMessageProcessor(opts)
	if err != nil {
		log.Error(err)
		return mp.Continue()
	}
	return mp.Process()
}
