package types

import "context"

const (
	ArgsKey           string = "args"
	CommandKey        string = "command"
	MidiKey           string = "midi"
	MidiBatchKey      string = "batch"
	MidiIdKey         string = "id"
	MidiMessagesKey   string = "messages"
	MidiDeviceKey     string = "device"
	MidiPitchKey      string = "pitch"
	MidiVelocityKey   string = "velocity"
	MidiNoteOffKey    string = "note_off"
	MidiNoteOnKey     string = "note_on"
	MidiCCKey         string = "cc"
	MidiControllerKey string = "controller"
	MidiValueKey      string = "value"
)

// CLI Flag types
type ParserName string
type ParserType ParserName
type Flags struct {
	Args        []string
	Command     CommandType
	Daemon      bool
	Example     bool
	ListDevices bool
	LogLevel    string
	Parser      ParserType
	Version     bool
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
type MidiOps map[MidiOpType]interface{}
type MidiArgs struct {
	Id      string
	Device  uint8
	Channel uint8
	NoteOn  MidiNoteOn
	NoteOff uint8
	CC      MidiCC
}

type MidiCall struct {
	Id   int
	Op   MidiOpType
	Args *MidiArgs
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

func MidiTempoType() MidiOpType {
	return MidiOpType("tempo_bpm")
}

func MidiCCType() MidiOpType {
	return MidiOpType("cc")
}
