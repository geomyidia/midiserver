package server

import "errors"

var (
	ErrUnsupMessageType = errors.New("unsupported message type")
)
