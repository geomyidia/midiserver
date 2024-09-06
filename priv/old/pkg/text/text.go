package text

import (
	log "github.com/sirupsen/logrus"

	"github.com/geomyidia/erlcmd/pkg/messages"
)

type Response struct {
	result messages.Result
	err    messages.Err
}

func NewResponse(result messages.Result, err messages.Err) *Response {
	return &Response{
		result: result,
		err:    err,
	}
}

func (r *Response) Send() {
	if r.err != "" {
		log.Error(r.err)
	} else {
		if r.result != "ok" {
			println(string(r.result))
		}
	}
}
