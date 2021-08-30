package messages

import (
	"os"

	erlang "github.com/okeuday/erlang_go/v2/erlang"
	log "github.com/sirupsen/logrus"

	"github.com/ut-proj/midiserver/pkg/types"
)

type Response struct {
	hasError bool
	result   erlang.OtpErlangTuple
	err      erlang.OtpErlangTuple
}

func NewResponse(result types.Result, err types.Err) *Response {
	hasError := false
	if err != "" {
		hasError = true
	}
	msg := &Response{
		result:   ResultMessage(result),
		err:      ErrorMessage(err),
		hasError: hasError,
	}
	log.Debugf("created result message: %#v", msg)
	return msg
}

// ErrorMessage ...
func ErrorMessage(errMsg types.Err) erlang.OtpErlangTuple {
	err := make([]interface{}, 2)
	err[0] = erlang.OtpErlangAtom("error")
	err[1] = errMsg
	log.Warnf("created error tuple: %#v", err)
	return erlang.OtpErlangTuple(err)
}

// ResultMessage ...
func ResultMessage(value types.Result) erlang.OtpErlangTuple {
	result := make([]interface{}, 2)
	result[0] = erlang.OtpErlangAtom("result")
	result[1] = value
	log.Tracef("created result tuple: %#v", result)
	return erlang.OtpErlangTuple(result)
}

// SendMessage ...
func (r *Response) Send() {
	msg := r.result
	if r.hasError {
		msg = r.err
		log.Errorf("Response: %+v", msg)

	}

	bytes, err := erlang.TermToBinary(msg, -1)
	if err != nil {
		log.Error(err)
		return
	}
	os.Stdout.Write(bytes)
	os.Stdout.Write([]byte("\n"))
}
