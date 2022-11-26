package messages

import "errors"

var (
	ErrMsgListFormat  = errors.New("too many elements in message list")
	ErrMsgNameFormat  = errors.New("incorrect message name format (needs to be atom)")
	ErrMsgTupleFormat = errors.New("incorrect message format (message needs to be tuple)")
	ErrMsgAtomFormat  = errors.New("incorrect message format (command needs to be atom)")
	ErrMsgValueFormat = errors.New("message value in unsupported format (needs to be atom , list, or tuple)")
)
