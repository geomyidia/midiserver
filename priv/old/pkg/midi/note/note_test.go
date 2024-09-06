package note

// Basic imports
import (
	"testing"

	"github.com/stretchr/testify/suite"
)

type MidiNotesTestSuite struct {
	suite.Suite
}

func (suite *MidiNotesTestSuite) SetupTest() {
}

func (suite *MidiNotesTestSuite) TestDiatonic() {
	//suite.NoError(err)
	suite.Equal(uint8(0), C)
	suite.Equal(uint8(2), D)
	suite.Equal(uint8(4), E)
	suite.Equal(uint8(5), F)
	suite.Equal(uint8(7), G)
	suite.Equal(uint8(9), A)
	suite.Equal(uint8(11), B)
}

func (suite *MidiNotesTestSuite) TestFlats() {
	suite.Equal(uint8(1), Db)
	suite.Equal(uint8(3), Eb)
	suite.Equal(uint8(6), Gb)
	suite.Equal(uint8(8), Ab)
	suite.Equal(uint8(10), Bb)
}

func (suite *MidiNotesTestSuite) TestSharps() {
	suite.Equal(uint8(1), Cs)
	suite.Equal(uint8(3), Ds)
	suite.Equal(uint8(6), Fs)
	suite.Equal(uint8(8), Gs)
	suite.Equal(uint8(10), As)
}

func (suite *MidiNotesTestSuite) TestDoubleFlats() {
	suite.Equal(uint8(0), Dbb)
	suite.Equal(uint8(2), Ebb)
	suite.Equal(uint8(5), Gbb)
	suite.Equal(uint8(7), Abb)
	suite.Equal(uint8(9), Bbb)
}

func (suite *MidiNotesTestSuite) TestDoubleSharps() {
	suite.Equal(uint8(2), Css)
	suite.Equal(uint8(4), Dss)
	suite.Equal(uint8(7), Fss)
	suite.Equal(uint8(9), Gss)
	suite.Equal(uint8(11), Ass)
}

// In order for 'go test' to run this suite, we need to create
// a normal test function and pass our suite to suite.Run
func TestMidiNotesTestSuite(t *testing.T) {
	suite.Run(t, new(MidiNotesTestSuite))
}
