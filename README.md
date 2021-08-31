# midiserver

[![Build Status][gh-actions-badge]][gh-actions]
[![Go Versions][go badge]][go]
[![LFE Versions][lfe badge]][lfe]
[![Erlang Versions][erlang badge]][erlang]

[![][logo]][logo-large]

*A MIDI CLI tool and server written in Go, focused on supporting BEAM music apps via Erlang Ports*

## Usage

```shell
$ ./bin/midiserver -h
```
```text
Usage: ./bin/midiserver [flags] [commands] [args]

Flags:

  -d    Daemonise midiserver; this disables the text parser; short-form flag
  -daemon
        Daemonise midiserver; this disables the text parser
  -l string
        Set the logging level; short-form flag (default "warn")
  -loglevel string
        Set the logging level (default "warn")
  -p string
        Set the parser to user for commands and data. Legal values are:
        [exec, port, text]. Note that setting to 'text' disables
        daemonisation and setting any of the other parsers automatically 
        enables daemonisation; short-form flag (default "text")
  -parser string
        Set the parser to user for commands and data. Legal values are:
        [exec, port, text]. Note that setting to 'text' disables
        daemonisation and setting any of the other parsers automatically 
        enables daemonisation (default "text")
  -v    Display version/build info and exit; short-form flag
  -version
        Display version/build info and exit

Commands:

  example [args]
        An example piece of music will be played on given device and channel.
        Valid device numbers are any of the "out" devices in the output of
        the 'list-devices' command; valid channel numbers are any of the 16
        MIDI channels: 0 through 15.
  list-devices
        will list the MIDI devices currently recognised by the operating
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
        provided for testing purposes by Erlang Ports implementations
  version
        an alternate form of the version info with concise formatting
```

## License

Apache Version 2 License

Copyright © 2020-2021, Duncan McGreggor

[//]: ---Named-Links---

[logo]: assets/images/logo-v1-x250.png
[logo-large]: assets/images/logo-v1-x1000.png
[github]: https://github.com/ut-proj/midiserver
[gh-actions-badge]: https://github.com/ut-proj/midiserver/workflows/ci%2Fcd/badge.svg
[gh-actions]: https://github.com/ut-proj/midiserver/actions
[go]: https://golang.org/
[go badge]: https://img.shields.io/badge/go-1.16-blue.svg
[lfe]: https://github.com/lfe/lfe
[lfe badge]: https://img.shields.io/badge/lfe-2.0-blue.svg
[erlang badge]: https://img.shields.io/badge/erlang-21%20to%2024-blue.svg
[erlang]: https://github.com/ut-proj/midiserver/blob/master/.github/workflows/cicd.yml
