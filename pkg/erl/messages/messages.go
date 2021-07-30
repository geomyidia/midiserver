package messages

import (
	"os"

	erlang "github.com/okeuday/erlang_go/v2/erlang"
	log "github.com/sirupsen/logrus"
)

// ErrorMessage ...
func ErrorMessage(errMsg string) erlang.OtpErlangTuple {
	err := make([]interface{}, 2)
	err[0] = erlang.OtpErlangAtom("error")
	err[1] = errMsg
	log.Warnf("Created error tuple: %#v", err)
	return erlang.OtpErlangTuple(err)
}

// ResultMessage ...
func ResultMessage(value string) erlang.OtpErlangTuple {
	result := make([]interface{}, 2)
	result[0] = erlang.OtpErlangAtom("result")
	result[1] = value
	log.Tracef("Created result tuple: %#v", result)
	return erlang.OtpErlangTuple(result)
}

// SendMessage ...
func SendMessage(tuple erlang.OtpErlangTuple) {
	msg, err := erlang.TermToBinary(tuple, -1)
	if err != nil {
		log.Error(err)
	}
	os.Stdout.Write(msg)
	os.Stdout.Write([]byte("\n"))
}

// SendError ...
func SendError(errMsg string) {
	log.Error(errMsg)
	err := ErrorMessage(errMsg)
	SendMessage(err)
}

// SendResult ...
func SendResult(value string) {
	msg := ResultMessage(value)
	log.Debugf("Created result message: %#v", msg)
	SendMessage(msg)
}
