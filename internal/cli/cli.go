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
        An example piece of music will be played on given device and channel.
        Valid device numbers are any of the "out" devices in the output of
        the 'list-devices' command; valid channel numbers are any of the 16
        MIDI channels: 0 through 15.
  list-devices
        This will list the MIDI devices currently recognised by the operating
        system, grouped by input devices and output devices.
  play-note [args]
        A pitch will be played with the default values for the arguments,
        opertionally overridden. Positional args are the integer values for
        device, channel, pitch, velocity, and duration.  Valid device numbers
        are any of the "out" devices in the output of the 'list-devices'
        command; valid channel numbers are any of the 16 MIDI channels:
        0 through 15. Pitch and velocity are standard MIDI integer values for
        the same. Duration is in seconds.
  ping
        Provided for testing purposes by Erlang Ports implementations.
  version
        An alternate form of the version info with concise formatting.

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

	var loglevel string
	logDef := "warn"
	logUse := "Set the logging level"
	flag.StringVar(&loglevel, "loglevel", logDef, logUse)
	flag.StringVar(&loglevel, "l", logDef, logUse+shortUse)

	var parser string
	parsDef := DefaultParser
	parsLegals := "[" +
		ExecParser + ", " +
		PortParser + ", " +
		TextParser + "]"
	parsUse := "Set the parser to user for commands and data. Legal values are:\n" +
		parsLegals + ". Note that setting to 'text' disables\n" +
		"daemonisation and setting any of the other parsers automatically \n" +
		"enables daemonisation"
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
		Args:     flag.Args(),
		Command:  types.Command(cmd),
		Daemon:   daemon,
		LogLevel: loglevel,
		Parser:   types.Parser(types.ParserName(parser)),
		Version:  vsn,
	}

	if flags.Version {
		println("Version: ", version.VersionString())
		println("Build: ", version.BuildString())
		fmt.Printf("Go: %s (%s)\n", version.GoVersionString(), version.GoArchString())
		os.Exit(0)
	}
	return flags
}
