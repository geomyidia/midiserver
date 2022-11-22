package datatypes

import "errors"

var (
	ErrCastingTuple  = errors.New("couldn't cast interface as Erlang tuple")
	ErrCastingList   = errors.New("couldn't cast interface as Erlang list")
	ErrCastingAtom   = errors.New("couldn't cast interface as Erlang atom")
	ErrBadTupleArity = errors.New("unexpected tuple arity (wrong number of elements")
)
