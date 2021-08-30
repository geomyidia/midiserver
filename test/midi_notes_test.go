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
	suite.Equal(0, note.C)
	suite.Equal(2, note.D)
	suite.Equal(4, note.E)
	suite.Equal(5, note.F)
	suite.Equal(7, note.G)
	suite.Equal(9, note.A)
	suite.Equal(11, note.B)
}

func (suite *MidiNotesTestSuite) TestFlats() {
	suite.Equal(1, note.Db)
	suite.Equal(3, note.Eb)
	suite.Equal(6, note.Gb)
	suite.Equal(8, note.Ab)
	suite.Equal(10, note.Bb)
}

func (suite *MidiNotesTestSuite) TestSharps() {
	suite.Equal(1, note.Cs)
	suite.Equal(3, note.Ds)
	suite.Equal(6, note.Fs)
	suite.Equal(8, note.Gs)
	suite.Equal(10, note.As)
}

func (suite *MidiNotesTestSuite) TestDoubleFlats() {
	suite.Equal(0, note.Dbb)
	suite.Equal(2, note.Ebb)
	suite.Equal(5, note.Gbb)
	suite.Equal(7, note.Abb)
	suite.Equal(9, note.Bbb)
}

func (suite *MidiNotesTestSuite) TestDoubleSharps() {
	suite.Equal(2, note.Css)
	suite.Equal(4, note.Dss)
	suite.Equal(7, note.Fss)
	suite.Equal(9, note.Gss)
	suite.Equal(11, note.Ass)
}

// In order for 'go test' to run this suite, we need to create
// a normal test function and pass our suite to suite.Run
func TestMidiNotesTestSuite(t *testing.T) {
	suite.Run(t, new(MidiNotesTestSuite))
}
