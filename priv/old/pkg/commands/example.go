package commands

import (
	log "github.com/sirupsen/logrus"

	"github.com/ut-proj/midiserver/pkg/midi"
	"github.com/ut-proj/midiserver/pkg/midi/note"
	"github.com/ut-proj/midiserver/pkg/types"
)

type PlayExampleOpts struct {
	DeviceId    uint8
	MidiChannel uint8
}

func DefaultPlayExampleOpts() *PlayExampleOpts {
	return &PlayExampleOpts{
		DeviceId:    0,
		MidiChannel: 0,
	}
}
func PlayExample(args map[string]interface{}) error {
	var opts *PlayExampleOpts
	if len(args) == 0 {
		log.Debug("got nil args ...")
		opts = DefaultPlayExampleOpts()
	} else {
		opts = &PlayExampleOpts{
			DeviceId:    args["device"].(uint8),
			MidiChannel: args["channel"].(uint8),
		}
	}
	log.Debugf("Got opts: %+v", opts)

	sys := midi.NewSystem()
	err := sys.SetWriter(opts.DeviceId)
	if err != nil {
		return err
	}
	err = sys.SetWriterChannel(opts.MidiChannel)
	if err != nil {
		return err
	}
	// A rendition of Max Richter's "On the Nature of Daylight" from the album
	// The Blue Notebooks
	chord1 := types.Chord(50, 3600, []uint8{note.Bb, 2}, []uint8{note.F, 3}, []uint8{note.Db, 4})
	chord2 := types.Chord(50, 3600, []uint8{note.Ab, 2}, []uint8{note.F, 3}, []uint8{note.C, 4})
	chord3 := types.Chord(50, 3600, []uint8{note.Db, 2}, []uint8{note.Db, 3}, []uint8{note.Ab, 3})
	chord4 := types.Chord(50, 3600, []uint8{note.Gb, 2}, []uint8{note.Gb, 3}, []uint8{note.Bb, 3})
	chord5 := types.Chord(50, 3600, []uint8{note.C, 3}, []uint8{note.Ab, 3}, []uint8{note.Eb, 4})
	chord6 := types.Chord(50, 3600, []uint8{note.Db, 3}, []uint8{note.Ab, 3}, []uint8{note.Eb, 4})
	chord7 := types.Chord(50, 3600, []uint8{note.C, 3}, []uint8{note.Ab, 3}, []uint8{note.F, 4})
	chord8 := types.Chord(50, 3600, []uint8{note.F, 2}, []uint8{note.Db, 3}, []uint8{note.Ab, 3})
	chord9 := types.Chord(50, 3600, []uint8{note.Ab, 3}, []uint8{note.F, 3}, []uint8{note.C, 4})
	chord10 := types.Chord(50, 3600, []uint8{note.Gb, 2}, []uint8{note.Eb, 3}, []uint8{note.Bb, 3})
	chord11 := types.Chord(50, 3600, []uint8{note.Eb, 2}, []uint8{note.Gb, 3}, []uint8{note.C, 4})
	chords := []*types.MidiChord{
		chord1, chord2, chord3, chord4,
		chord1, chord2, chord3, chord4,
		chord1, chord5, chord6, chord7,
		chord1, chord5, chord6, chord2,
		chord1, chord8, chord4, chord9,
		chord1, chord8, chord10, chord11,
	}
	for _, chord := range chords {
		err = sys.PlayChord(chord)
		if err != nil {
			return nil
		}
	}

	return nil
}
