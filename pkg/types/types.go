package types

import "context"

const (
	ArgsKey string = "args"
	CommandKey string = "command"
)

type CommandName string
type CommandType CommandName
type ParserName string
type ParserType ParserName
type Result string
type Err string
type Flags struct {
	Args     []string
	Command  CommandType
	Daemon   bool
	LogLevel string
	Parser   ParserType
	Version  bool
}
type Proplist map[string]interface{}
type CommandProcessor func(context.Context, CommandType, Proplist, *Flags)
type MessageProcessor func() Result

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

func Command(name CommandName) CommandType {
	return CommandType(name)
}

func ExampleCommand() CommandType {
	return CommandType(CommandType("example"))
}

func ListDevicesCommand() CommandType {
	return CommandType(CommandType("midi"))
}

func MidiCommand() CommandType {
	return CommandType(CommandType("midi"))
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

func (r Result) ToCommand() CommandType {
	return Command(CommandName(string(r)))
}