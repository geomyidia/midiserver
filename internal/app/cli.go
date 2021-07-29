package app

import (
	"flag"
	"fmt"
	"os"
)

type Flags struct {
	LogLevel string
	Version  bool
}

func ParseCLI() *Flags {
	logLevelPtr := flag.String("loglevel", "info", "Set the logging level")
	versionPtr := flag.Bool("version", false, "Display version/build info and exit")
	flag.Parse()
	flags := &Flags{
		LogLevel: *logLevelPtr,
		Version:  *versionPtr,
	}

	if flags.Version {
		println("Version: ", VersionString())
		println("Build: ", BuildString())
		fmt.Printf("Go: %s (%s)\n", GoVersionString(), GoArchString())
		os.Exit(0)
	}
	return flags
}
