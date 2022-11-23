package datatypes

import "errors"

var (
	ErrBadTupleArity      = errors.New("unexpected tuple arity (wrong number of elements")
	ErrCastingAtom        = errors.New("couldn't cast interface as Erlang atom")
	ErrCastingList        = errors.New("couldn't cast interface as Erlang list")
	ErrCastingTuple       = errors.New("couldn't cast interface as Erlang tuple")
	ErrUnsupportedGoType  = errors.New("unsupported Go type")
	ErrUnsupportedOTPType = errors.New("unsupported Erlang/OTP type")
	ErrNotImplemented     = errors.New("feature not yet implemented")
)
