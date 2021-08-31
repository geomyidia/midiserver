package midi

// Basic imports
import (
	"testing"

	"github.com/stretchr/testify/suite"

	"github.com/ut-proj/midiserver/pkg/midi/note"
)

type MidiNotesTestSuite struct {
	suite.Suite
}

func (suite *MidiNotesTestSuite) SetupTest() {
}

func (suite *MidiNotesTestSuite) TestDiatonic() {
	//suite.NoError(err)
	suite.Equal(uint8(0), note.C)
	suite.Equal(uint8(2), note.D)
	suite.Equal(uint8(4), note.E)
	suite.Equal(uint8(5), note.F)
	suite.Equal(uint8(7), note.G)
	suite.Equal(uint8(9), note.A)
	suite.Equal(uint8(11), note.B)
}

func (suite *MidiNotesTestSuite) TestFlats() {
	suite.Equal(uint8(1), note.Db)
	suite.Equal(uint8(3), note.Eb)
	suite.Equal(uint8(6), note.Gb)
	suite.Equal(uint8(8), note.Ab)
	suite.Equal(uint8(10), note.Bb)
}

func (suite *MidiNotesTestSuite) TestSharps() {
	suite.Equal(uint8(1), note.Cs)
	suite.Equal(uint8(3), note.Ds)
	suite.Equal(uint8(6), note.Fs)
	suite.Equal(uint8(8), note.Gs)
	suite.Equal(uint8(10), note.As)
}

func (suite *MidiNotesTestSuite) TestDoubleFlats() {
	suite.Equal(uint8(0), note.Dbb)
	suite.Equal(uint8(2), note.Ebb)
	suite.Equal(uint8(5), note.Gbb)
	suite.Equal(uint8(7), note.Abb)
	suite.Equal(uint8(9), note.Bbb)
}

func (suite *MidiNotesTestSuite) TestDoubleSharps() {
	suite.Equal(uint8(2), note.Css)
	suite.Equal(uint8(4), note.Dss)
	suite.Equal(uint8(7), note.Fss)
	suite.Equal(uint8(9), note.Gss)
	suite.Equal(uint8(11), note.Ass)
}

// In order for 'go test' to run this suite, we need to create
// a normal test function and pass our suite to suite.Run
func TestMidiNotesTestSuite(t *testing.T) {
	suite.Run(t, new(MidiNotesTestSuite))
}
