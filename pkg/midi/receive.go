package midi

import (
	log "github.com/sirupsen/logrus"
	"gitlab.com/gomidi/midi"
	"gitlab.com/gomidi/midi/reader"
)

func ReceiveEach(p *reader.Position, msg midi.Message) {
	log.Debugf("got MIDI msg %+v (at position %v)", msg, p)
}
