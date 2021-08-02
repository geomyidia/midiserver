package types

import "context"

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
type CommandProcessor func(context.Context, CommandType, *Flags)
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
	return CommandType(ParserName("ping"))
}

func ListDevicesCommand() CommandType {
	return CommandType(ParserName("midi"))
}

func MidiCommand() CommandType {
	return CommandType(ParserName("midi"))
}

func PingCommand() CommandType {
	return CommandType(ParserName("ping"))
}

func StopCommand() CommandType {
	return CommandType(ParserName("stop"))
}

func VersionCommand() CommandType {
	return CommandType(CommandType("version"))
}
