package midi

import (
	log "github.com/sirupsen/logrus"
	"gitlab.com/gomidi/midi"
	"gitlab.com/gomidi/midi/reader"

	"github.com/ut-proj/midiserver/pkg/erl/rpc"
)

func ReceiveEach(p *reader.Position, msg midi.Message) {
	log.Debugf("got MIDI msg %+v (at position %v)", msg, p)
}

func ReceiveUnknown(p *reader.Position, msg midi.Message) {
	log.Debugf("got unknown msg %+v (at position %v)", msg, p)
}

func ReceiveClock(rpcClient *rpc.Client) func() {
	return func() {
		val, err := rpcClient.ClockTick()
		if err != nil {
			log.Error(err)
		}
		if val != "" {
			log.Trace(val)
		}
		log.Trace("got clock msg")
	}
}

func ReceiveContinue() {
	log.Debug("got continue msg")
}

func ReceiveReset() {
	log.Debug("got reset msg")
}

func ReceiveStart() {
	log.Debug("got start msg")
}

func ReceiveStop() {
	log.Debug("got stop msg")
}

func ReceiveTick() {
	log.Debug("got tick msg")
}
