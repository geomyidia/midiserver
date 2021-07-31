package port

import (
	log "github.com/sirupsen/logrus"

	"github.com/geomyidia/erl-midi-server/pkg/erl/term"
)

// ProcessMessage ...
func ProcessMessage() term.Result {
	opts := term.DefaultOpts()
	mp, err := term.NewMessageProcessor(opts)
	if err != nil {
		log.Error(err)
		return mp.Continue()
	}
	return mp.Process()
}
