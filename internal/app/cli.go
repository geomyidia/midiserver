package app

import (
	"flag"
	"fmt"
	"os"
)

const (
	DefaultParser string = "text"
	ExecParser    string = "exec"
	PortParser    string = "port"
	TextParser    string = "text"
)

type Flags struct {
	Daemon   bool
	LogLevel string
	Parser   string
	Version  bool
}

func ParseCLI() *Flags {
	shortUse := "; short-form flag"

	var daemon bool
	daemDef := true
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
		parsLegals
	flag.StringVar(&parser, "parser", parsDef, parsUse)
	flag.StringVar(&parser, "p", parsDef, parsUse+shortUse)

	var version bool
	verDef := false
	verUse := "Display version/build info and exit"
	flag.BoolVar(&version, "version", verDef, verUse)
	flag.BoolVar(&version, "v", verDef, verUse+shortUse)

	flag.Parse()
	flags := &Flags{
		Daemon:   daemon,
		LogLevel: loglevel,
		Parser:   parser,
		Version:  version,
	}

	if flags.Version {
		println("Version: ", VersionString())
		println("Build: ", BuildString())
		fmt.Printf("Go: %s (%s)\n", GoVersionString(), GoArchString())
		os.Exit(0)
	}
	return flags
}
