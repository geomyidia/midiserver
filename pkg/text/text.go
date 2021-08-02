package text

import (
	log "github.com/sirupsen/logrus"

	"github.com/geomyidia/erl-midi-server/pkg/types"
)

type Response struct {
	result types.Result
	err    types.Err
}

func NewResponse(result types.Result, err types.Err) *Response {
	return &Response{
		result: result,
		err:    err,
	}
}

func (r *Response) Send() {
	if r.err != "" {
		log.Error(r.err)
	} else {
		println(string(r.result))
	}
}
