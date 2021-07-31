package erl

import (
	"bufio"
	"encoding/hex"
	"errors"
	"fmt"
	"os"

	erlang "github.com/okeuday/erlang_go/v2/erlang"
	log "github.com/sirupsen/logrus"

	"github.com/geomyidia/erl-midi-server/pkg/erl/messages"
)

// Constants
const (
	DELIMITER       = '\n'
	DRCTVARITY      = 2
	DRCTVKEYINDEX   = 0
	DRCTVVALUEINDEX = 1
)

type Result string

type Opts struct {
	IsHexEncoded bool
}

func DefaultOpts() *Opts {
	return &Opts{
		IsHexEncoded: false,
	}
}

type Packet struct {
	bytes     []byte
	len       int
	last      int
	trimmed   []byte
	unwrapped []byte
	opts      *Opts
}

func ReadStdIOPacket(opts *Opts) (*Packet, error) {
	reader := bufio.NewReader(os.Stdin)
	bytes, _ := reader.ReadBytes(DELIMITER)
	byteLen := len(bytes)
	if byteLen == 0 {
		return nil, errors.New("read zero bytes")
	}
	log.Tracef("Original packet: %#v", bytes)
	log.Tracef("Original packet length: %d", byteLen)
	packet := &Packet{
		bytes: bytes,
		len:   byteLen,
		last:  byteLen - 1,
	}
	packet.setTrimmed()
	err := packet.setUnwrapped()
	if err != nil {
		return nil, err
	}
	return packet, nil
}

func (p *Packet) setTrimmed() {
	p.trimmed = p.bytes[:p.last]
}

func (p *Packet) Bytes() []byte {
	if p.opts.IsHexEncoded {
		return p.unwrapped
	}
	return p.trimmed
}

// setUnwrapped is a utility method for a hack needed in order to
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
func (p *Packet) setUnwrapped() error {
	if p.opts.IsHexEncoded {
		hexStr := string(p.trimmed[:])
		bytes, err := hex.DecodeString(hexStr)
		if err != nil {
			return fmt.Errorf("problem unwrapping packet: %s", err.Error())
		}
		p.trimmed = bytes
	}
	return nil
}

func (p *Packet) Term() (interface{}, error) {
	bytes := p.Bytes()
	term, err := erlang.BinaryToTerm(bytes)
	if err != nil {
		return nil, fmt.Errorf("problem creating Erlang term from %#v: %s",
			bytes, err.Error())
	}
	return term, nil
}

type Message struct {
	tuple     erlang.OtpErlangTuple
	directive erlang.OtpErlangAtom
	payload   interface{}
}

func NewMessage(t interface{}) (*Message, error) {
	tuple, ok := t.(erlang.OtpErlangTuple)
	if !ok {
		return nil, errors.New("unexpected message format")
	}
	if len(tuple) != DRCTVARITY {
		return nil, fmt.Errorf("tuple of wrong size; expected 2, got %d", len(tuple))
	}
	directive, ok := tuple[DRCTVKEYINDEX].(erlang.OtpErlangAtom)
	if !ok {
		return nil, errors.New("unexpected type for directive")
	}
	msg := &Message{tuple: tuple}
	msg.directive = directive
	msg.payload = tuple[DRCTVVALUEINDEX]
	return msg, nil
}

func (m *Message) Directive() erlang.OtpErlangAtom {
	return m.directive
}

func (m *Message) Payload() interface{} {
	return m.payload
}

func (m *Message) IsCommand() bool {
	return m.directive == erlang.OtpErlangAtom("command")
}

func (m *Message) IsMIDI() bool {
	return m.directive == erlang.OtpErlangAtom("midi")
}

func (m *Message) Command() (erlang.OtpErlangAtom, error) {
	if !m.IsCommand() {
		return erlang.OtpErlangAtom("error"),
			errors.New("directive is not a command")
	}
	command, ok := m.Payload().(erlang.OtpErlangAtom)
	if !ok {
		return erlang.OtpErlangAtom("error"),
			errors.New("could not extract command atom")
	}
	return command, nil
}

type MessageProcessor struct {
	packet *Packet
	term   interface{}
	msg    *Message
}

func NewMessageProcessor(opts *Opts) (*MessageProcessor, error) {
	packet, err := ReadStdIOPacket(opts)
	if err != nil {
		return &MessageProcessor{}, err
	}
	t, err := packet.Term()
	if err != nil {
		return &MessageProcessor{}, err
	}
	log.Debugf("Got Erlang Port term")
	log.Tracef("%#v", t)
	msg, err := NewMessage(t)
	if err != nil {
		messages.SendError(err.Error())
		return &MessageProcessor{}, err
	}
	return &MessageProcessor{
		packet: packet,
		term:   t,
		msg:    msg,
	}, nil
}

func (mp *MessageProcessor) Continue() Result {
	return Result("continue")
}

func (mp *MessageProcessor) Process() Result {
	if mp.msg.IsCommand() {
		command, err := mp.msg.Command()
		if err != nil {
			log.Error(err)
			return mp.Continue()
		}
		return Result(command)
	} else if mp.msg.IsMIDI() {
		// process MIDI message
		return mp.Continue()
	} else {
		log.Error("Unexected message type")
		return mp.Continue()
	}
}
