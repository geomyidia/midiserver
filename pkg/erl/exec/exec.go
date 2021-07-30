package exec

import (
	"bufio"
	"encoding/hex"
	"fmt"
	"os"

	erlang "github.com/okeuday/erlang_go/v2/erlang"
	log "github.com/sirupsen/logrus"

	"github.com/geomyidia/erl-midi-server/pkg/erl/messages"
)

// Constants
const (
	READERSIZE    = 1677216
	DELIMITER     = '\n'
	CMDARITY      = 2
	CMDKEYINDEX   = 0
	CMDVALUEINDEX = 1
)

// ProcessMessage - when using the Erlang exec library to send
// messages to this Go server, a stange thing happens: a byte is
// dropped from the middle of the
func ProcessMessage() string {
	reader := bufio.NewReader(os.Stdin)
	packet, _ := reader.ReadBytes(DELIMITER)
	if len(packet) == 0 {
		log.Fatal("Read zero bytes")
		return "continue"
	}
	log.Tracef("Original packet: %#v", packet)
	log.Tracef("Original packet length: %d", len(packet))
	unwrapped, err := unwrap(packet)
	if err != nil {
		log.Errorf("Problem unwrapping packet ...")
		log.Error(err)
		return "continue"
	}
	log.Tracef("Unwrapped packet: %#v", unwrapped)
	term, err := erlang.BinaryToTerm(unwrapped)
	if err != nil {
		log.Errorf("Problem with packet: %#v", packet)
		log.Error(err)
		return "continue"
	}
	log.Debugf("Got Erlang Exec message: %#v", term)
	tuple, ok := term.(erlang.OtpErlangTuple)
	if !ok {
		messages.SendError("Did not receive expected message type")
		return "continue"
	}
	if len(tuple) != CMDARITY {
		messages.SendError(fmt.Sprintf("Tuple of wrong size; expected 2, got %d", len(tuple)))
		return "continue"
	}
	if tuple[CMDKEYINDEX] != erlang.OtpErlangAtom("command") {
		messages.SendError("Did not receive expected tuple message format {command, ...}")
		return "continue"
	}
	command, ok := tuple[CMDVALUEINDEX].(erlang.OtpErlangAtom)
	if !ok {
		messages.SendError("Did not receive command as Erlang atom")
		return "continue"
	}
	return string(command)
}

// unwrap is a utility function for a hack needed in order to
// successully handle messages from the Erlang exec library.
//
// What was happening when exec messages were being processed
// by ProcessPortMessage was that a single byte was being dropped
// from the middle (in the case of the #(command ping) message,
// it was byte 0x04 of the Term protocol encoded bytes). The
// bytes at the sending end were present and correct, just not
// at the receiving end.
//
// So, in order to get around this, the sending end hex-encoded
// the Term protocol bytes and send that as a bitstring; the
// function below hex-decodes this, and allows the function
// ProcessExecMessage to handle binary encoded Term data with
// none of its bytes missing.
func unwrap(hexBytes []byte) ([]byte, error) {
	trimedHexBytes := hexBytes[:len(hexBytes)-1]
	hexStr := string(trimedHexBytes[:])
	return hex.DecodeString(hexStr)
}
