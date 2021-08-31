package midi

import (
	"time"

	log "github.com/sirupsen/logrus"
	"gitlab.com/gomidi/midi/writer"

	"github.com/ut-proj/midiserver/pkg/types"
)

func (s *System) PlayChord(chord *types.MidiChord) error {
	log.Debug("playing chord ...")
	for _, pitch := range chord.Pitches {
		err := writer.NoteOn(s.Writer, pitch, chord.Velocity)
		log.Tracef("\tpitch: %v", pitch)
		if err != nil {
			return err
		}
	}
	time.Sleep(chord.Duration * time.Millisecond)
	for _, pitch := range chord.Pitches {
		err := writer.NoteOff(s.Writer, pitch)
		if err != nil {
			return err
		}
	}
	return nil
}
