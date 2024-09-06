// This package is really just meant to be used for playing some example notes. This application is
// not intended for composition, rather to support composition from Erlang/LFE/BEAM languages,
// sent to this application via Port messages. I guess you could use it to compose, if you wanted to?
// if you need more than what's here, let's chat, and submit a PR.
package note

// Chromatic scale with flats
const (
	C uint8 = iota
	Db
	D
	Eb
	E
	F
	Gb
	G
	Ab
	A
	Bb
	B
)

// Sharps
const (
	_ uint8 = iota
	Cs
	_
	Ds
	_
	_
	Fs
	_
	Gs
	_
	As
)

// Double flats
const (
	Dbb uint8 = iota
	_
	Ebb
	_
	_
	Gbb
	_
	Abb
	_
	Bbb
)

// Double sharps
const (
	_ uint8 = iota
	_
	Css
	_
	Dss
	_
	_
	Fss
	_
	Gss
	_
	Ass
)

func Pitches(notes ...[]uint8) []uint8 {
	var pitches []uint8
	for _, note := range notes {
		pitches = append(pitches, note[0]+(12*(1+note[1])))
	}
	return pitches
}
