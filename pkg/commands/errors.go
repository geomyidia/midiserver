package commands

import "errors"

var (
	ErrCmdMsgFormat = errors.New("incorrect message format (message type needs to be 'command')")
)
