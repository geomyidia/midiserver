package types

import (
	"context"
	"time"
)

const (
	ArgsKey              string = "args"
	CommandKey           string = "command"
	MidiKey              string = "midi"
	MidiBatchKey         string = "batch"
	MidiChannelKey       string = "channel"
	MidiIdKey            string = "id"
	MidiParallelKey      string = "parallel?"
	MidiMessagesKey      string = "messages"
	MidiDeviceKey        string = "device"
	MidiPitchKey         string = "pitch"
	MidiBankSelectMSBKey string = "bank_select_msb"
	MidiBankSelectLSBKey string = "bank_select_lsb"
	MidiProgramChangeKey string = "program_change"
	MidiVelocityKey      string = "velocity"
	MidiNoteOffKey       string = "note_off"
	MidiNoteOnKey        string = "note_on"
	MidiCCKey            string = "cc"
	MidiControllerKey    string = "controller"
	MidiValueKey         string = "value"
	MidiRealtimeKey      string = "realtime"
)

// CLI Flag types
type ParserName string
type ParserType ParserName
type Flags struct {
	Args            []string
	Command         CommandType
	Daemon          bool
	Example         bool
	ListDevices     bool
	LogLevel        string
	LogReportCaller bool
	MidiInDeviceID  int
	Parser          ParserType
	EPMDHost        string
	EPMDPort        int
	RemoteNode      string
	Version         bool
}

// General result types
type Result string
type Err string

// Command types
type CommandName string
type CommandType CommandName
type CommandProcessor func(context.Context, CommandType, PropList, *Flags)
type MessageProcessor func() Result

// MIDI types
type MidiOpType string
type MidiRTType string
type MidiPitch uint8
type MidiVelocity uint8
type MidiNoteOn struct {
	Pitch    uint8
	Velocity uint8
}
type MidiCC struct {
	Controller uint8
	Value      uint8
}
type MidiChord struct {
	Pitches  []uint8
	Velocity uint8
	Duration time.Duration
}
type MidiOps map[MidiOpType]interface{}
type MidiArgs struct {
	Id       string
	Device   uint8
	Channel  uint8
	NoteOn   MidiNoteOn
	NoteOff  uint8
	CC       MidiCC
	Chord    MidiChord
	Program  uint8
	Realtime MidiRTType
}

type MidiCall struct {
	Id   int
	Op   MidiOpType
	Args *MidiArgs
}

// Other types
type PropList map[string]interface{}

func Chord(velocity uint8, duration uint32, notes ...[]uint8) *MidiChord {
	var pitches []uint8
	for _, note := range notes {
		pitches = append(pitches, note[0]+(12*(1+note[1])))
	}
	return &MidiChord{
		Pitches:  pitches,
		Velocity: velocity,
		Duration: time.Duration(duration),
	}
}

// Part of CLI Options

func Parser(key ParserName) ParserType {
	return ParserType(key)
}

func ExecParser() ParserType {
	return ParserType(ParserName("exec"))
}

func PortParser() ParserType {
	return ParserType(ParserName("port"))
}

func TextParser() ParserType {
	return ParserType(ParserName("text"))
}

// Commands

func Command(name CommandName) CommandType {
	return CommandType(name)
}

func PlayNoteCommand() CommandType {
	return Command("play-note")
}

func ExampleCommand() CommandType {
	return Command("example")
}

func ListDevicesCommand() CommandType {
	return Command("list-devices")
}

func PingCommand() CommandType {
	return Command("ping")
}

func RemotePortCommand() CommandType {
	return Command("remote-port")
}

func StopCommand() CommandType {
	return Command("stop")
}

func VersionCommand() CommandType {
	return Command("version")
}

func EmptyCommand() CommandType {
	return Command("")
}

func (r Result) ToCommand() CommandType {
	return Command(CommandName(string(r)))
}

// MIDI

func MidiOp(name string) MidiOpType {
	return MidiOpType(name)
}

func MidiBatchType() MidiOpType {
	return MidiOpType("batch")
}

func MidiChannelType() MidiOpType {
	return MidiOpType("channel")
}

func MidiDeviceType() MidiOpType {
	return MidiOpType("device")
}

func MidiMeterType() MidiOpType {
	return MidiOpType("meter")
}

func MidiNoteOnType() MidiOpType {
	return MidiOpType("note_on")
}

func MidiNoteOffType() MidiOpType {
	return MidiOpType("note_off")
}

func MidiProgramChangeType() MidiOpType {
	return MidiOpType("program_change")
}

func MidiBankSelectMSBType() MidiOpType {
	return MidiOpType("bank_select_msb")
}

func MidiBankSelectLSBType() MidiOpType {
	return MidiOpType("bank_select_lsb")
}

func MidiTempoType() MidiOpType {
	return MidiOpType("tempo_bpm")
}

func MidiCCType() MidiOpType {
	return MidiOpType("cc")
}

func MidiChordType() MidiOpType {
	return MidiOpType("chord")
}

func MidiRealtimeType() MidiOpType {
	return MidiOpType("realtime")
}

func MidiRTClock() MidiRTType {
	return MidiRTType("clock")
}

func MidiRTContinue() MidiRTType {
	return MidiRTType("continue")
}

func MidiRTReset() MidiRTType {
	return MidiRTType("reset")
}

func MidiRTStart() MidiRTType {
	return MidiRTType("start")
}

func MidiRTStop() MidiRTType {
	return MidiRTType("stop")
}

func MidiRTTick() MidiRTType {
	return MidiRTType("tick")
}
