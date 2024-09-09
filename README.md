# midiserver

[![Build Status][gh-actions-badge]][gh-actions]
[![LFE Versions][lfe badge]][lfe]
[![Erlang Versions][erlang badge]][erlang]

[![][logo]][logo-large]

*A MIDI Erlang NIF for LFE*

## About

TBD

## Supported Versions

Dependencies:
* cmake
* make/gcc/clang

Currently this library is being developed using the following:

* macos 12 & 14
  * Erlang 26.2
  * rebar 3.24
* Raspberrian/Debian 12 (bookworm)
  * Erlang 25.2
  * rebar 3.19

Note that the NIF does not run properly under Erlang 27.

To run on a Raspberry PI (or other Debian distro), you will probably need the following:

```shell
sudo apt-get install cmake libasound2-dev
```

## Usage

TBD

## License

Apache Version 2 License

Copyright © 2020-2024, Duncan McGreggor

[//]: ---Named-Links---

[logo]: priv/images/logo-v1-x250.png
[logo-large]: priv/images/logo-v1-x1000.png
[gh-actions-badge]: https://github.com/ut-proj/midiserver/workflows/ci%2Fcd/badge.svg
[gh-actions]: https://github.com/ut-proj/midiserver/actions
[go]: https://golang.org/
[go badge]: https://img.shields.io/badge/go-1.16-blue.svg
[lfe]: https://github.com/lfe/lfe
[lfe badge]: https://img.shields.io/badge/lfe-2.0-blue.svg
[erlang badge]: https://img.shields.io/badge/erlang-21%20to%2024-blue.svg
[erlang]: https://github.com/ut-proj/midiserver/blob/master/.github/workflows/cicd.yml
