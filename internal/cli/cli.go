package cli

import (
	"flag"
	"fmt"
	"os"

	"github.com/geomyidia/erl-midi-server/pkg/version"
)

const (
	DefaultParser string = "text"
	ExecParser    string = "exec"
	PortParser    string = "port"
	TextParser    string = "text"
)

type Flags struct {
	Args     []string
	Daemon   bool
	LogLevel string
	Parser   string
	Version  bool
}

func Parse() *Flags {
	shortUse := "; short-form flag"

	var daemon bool
	daemDef := false
	daemUse := "Daemonise midiserver; this disables the text parser"
	flag.BoolVar(&daemon, "daemon", daemDef, daemUse)
	flag.BoolVar(&daemon, "d", daemDef, daemUse+shortUse)

	var loglevel string
	logDef := "info"
	logUse := "Set the logging level"
	flag.StringVar(&loglevel, "loglevel", logDef, logUse)
	flag.StringVar(&loglevel, "l", logDef, logUse+shortUse)

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
	flags := &Flags{
		Args:     flag.Args(),
		Daemon:   daemon,
		LogLevel: loglevel,
		Parser:   parser,
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
