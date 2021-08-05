package types

import "context"

const (
	ArgsKey    string = "args"
	CommandKey string = "command"
	MidiKey    string = "midi"
)

// CLI Flag types
type ParserName string
type ParserType ParserName
type Flags struct {
	Args     []string
	Command  CommandType
	Daemon   bool
	LogLevel string
	Parser   ParserType
	Version  bool
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
type MidiMeter struct {
	Numerator   uint8
	Denominator uint8
}
type MidiPitch uint8
type MidiVelocity uint8
type MidiNoteOn struct {
	Pitch    MidiPitch
	Velocity MidiVelocity
}
type MidiOps map[MidiOpType]interface{}
type MidiOpts struct {
	Device  uint8
	Meter   MidiMeter
	Tempo   uint8
	NoteOn  MidiNoteOn
	NoteOff uint8
}

// Other types
type PropList map[string]interface{}

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

func ExampleCommand() CommandType {
	return CommandType(CommandType("example"))
}

func ListDevicesCommand() CommandType {
	return CommandType(CommandType("list-devices"))
}

func PingCommand() CommandType {
	return CommandType(CommandType("ping"))
}

func StopCommand() CommandType {
	return CommandType(CommandType("stop"))
}

func VersionCommand() CommandType {
	return CommandType(CommandType("version"))
}

func EmptyCommand() CommandType {
	return CommandType(CommandType(""))
}

func (r Result) ToCommand() CommandType {
	return Command(CommandName(string(r)))
}

// MIDI

func MidiOp(name string) MidiOpType {
	return MidiOpType(name)
}

func MidiBatch() MidiOpType {
	return MidiOpType("batch")
}
