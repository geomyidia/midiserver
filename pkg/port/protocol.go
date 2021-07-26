package port

import (
	"bufio"
	"fmt"
	"os"

	erlang "github.com/okeuday/erlang_go/v2/erlang"
	log "github.com/sirupsen/logrus"
)

// Constants
const (
	READERSIZE    = 1677216
	DELIMITER     = '\n'
	CMDARITY      = 2
	CMDKEYINDEX   = 0
	CMDVALUEINDEX = 1
)

// ProcessPortMessage ...
func ProcessPortMessage() string {
	var term interface{}
	reader := bufio.NewReaderSize(os.Stdin, READERSIZE)
	packet, _ := reader.ReadBytes(DELIMITER)
	if len(packet) == 0 {
		log.Fatal("Read zero bytes")
		return "continue"
	}
	log.Debugf("Original packet: %#v", packet)
	log.Debugf("Packet length: %d", len(packet))
	// Drop the message's field separator, \0xa (newline)
	packet = packet[:len(packet)-1]
	log.Debugf("New packet: %#v", packet)
	term, err := erlang.BinaryToTerm(packet)
	if err != nil {
		log.Errorf("Problem with packet: %#v", packet)
		log.Error(err)
		return "continue"
	}
	log.Debugf("Got Erlang Port message: %#v", term)
	tuple, ok := term.(erlang.OtpErlangTuple)
	if !ok {
		SendError("Did not receive expected message type")
		return "continue"
	}
	if len(tuple) != CMDARITY {
		SendError(fmt.Sprintf("Tuple of wrong size; expected 2, got %d", len(tuple)))
		return "continue"
	}
	if tuple[CMDKEYINDEX] != erlang.OtpErlangAtom("command") {
		SendError("Did not receive expected tuple message format {command, ...}")
		return "continue"
	}
	command, ok := tuple[CMDVALUEINDEX].(erlang.OtpErlangAtom)
	if !ok {
		SendError("Did not receive command as Erlang atom")
		return "continue"
	}
	return string(command)
}
