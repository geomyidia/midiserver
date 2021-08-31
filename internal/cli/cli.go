package cli

import (
	"flag"
	"fmt"
	"os"

	"github.com/ut-proj/midiserver/pkg/types"
	"github.com/ut-proj/midiserver/pkg/version"
)

const (
	DefaultParser string = "text"
	ExecParser    string = "exec"
	PortParser    string = "port"
	TextParser    string = "text"
	helpText      string = `
  example [args]
        where the optional positional integer args are device and channel.
        An example piece of music will be played on given device and channel.
  list-devices
        will list the MIDI devices currently recognised by the operating system,
       grouped by input devices and output devices.
  play-note [args]
       where the optional positional integer args are device, channel,
       pitch, velocity, and duration. The pitch will be played with the
       given argument values for the given duration.
  ping
        provided for testing purposes by Erlang Ports implementations
  version
        an alternate form of the version info with concise formatting

	`
)

var (
	Usage = func() {
		w := flag.CommandLine.Output()
		fmt.Fprintf(w, "\nUsage: %s [flags] [commands] [args]\n", os.Args[0])
		fmt.Fprintf(w, "\nFlags:\n\n")
		flag.PrintDefaults()
		fmt.Fprintf(w, "\nCommands:\n%s", helpText)
	}
)

func Parse() *types.Flags {
	flag.Usage = Usage
	shortUse := "; short-form flag"

	var daemon bool
	daemDef := false
	daemUse := "Daemonise midiserver; this disables the text parser"
	flag.BoolVar(&daemon, "daemon", daemDef, daemUse)
	flag.BoolVar(&daemon, "d", daemDef, daemUse+shortUse)

	var example bool
	exampleDef := false
	exampleUse := "Play some example MIDI; required args are device, " +
		"channel, pitch, velocity, and duration"
	flag.BoolVar(&example, "example", exampleDef, exampleUse)

	var loglevel string
	logDef := "warn"
	logUse := "Set the logging level"
	flag.StringVar(&loglevel, "loglevel", logDef, logUse)
	flag.StringVar(&loglevel, "l", logDef, logUse+shortUse)

	var listing bool
	listDef := false
	listUse := "List the system MIDI devices"
	flag.BoolVar(&listing, "list", listDef, listUse)

	var parser string
	parsDef := DefaultParser
	parsLegals := "[" +
		ExecParser + ", " +
		PortParser + ", " +
		TextParser + "]"
	parsUse := "Set the parser to user for commands and data. Legal values are " +
		parsLegals + ". \n" +
		"Note that setting to 'text' disables daemonisation and setting any of " +
		"the other \n" +
		"parsers automatically enables daemonisation"
	flag.StringVar(&parser, "parser", parsDef, parsUse)
	flag.StringVar(&parser, "p", parsDef, parsUse+shortUse)

	var vsn bool
	verDef := false
	verUse := "Display version/build info and exit"
	flag.BoolVar(&vsn, "version", verDef, verUse)
	flag.BoolVar(&vsn, "v", verDef, verUse+shortUse)

	flag.Parse()
	var cmd types.CommandName
	args := flag.Args()
	if len(args) > 0 {
		cmd = types.CommandName(args[0])
	}
	flags := &types.Flags{
		Args:        flag.Args(),
		Command:     types.Command(cmd),
		Daemon:      daemon,
		ListDevices: listing,
		LogLevel:    loglevel,
		Parser:      types.Parser(types.ParserName(parser)),
		Version:     vsn,
	}

	if flags.Version {
		println("Version: ", version.VersionString())
		println("Build: ", version.BuildString())
		fmt.Printf("Go: %s (%s)\n", version.GoVersionString(), version.GoArchString())
		os.Exit(0)
	}
	return flags
}
