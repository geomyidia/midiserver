package messages

import "errors"

var (
	ErrCmdListFormat  = errors.New("too many elements in command list")
	ErrCmdTupleFormat = errors.New("incorrect command format (message needs to be tuple)")
	ErrCmdAtomFormat  = errors.New("incorrect command format (command needs to be atom)")
	ErrCmdValueFormat = errors.New("command value in unsupported format (needs to be atom or list)")
)
