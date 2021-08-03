package proplists

import (
	"errors"
	"fmt"

	erlang "github.com/okeuday/erlang_go/v2/erlang"
	log "github.com/sirupsen/logrus"

	"github.com/geomyidia/midiserver/pkg/types"
)

const (
	TUPLEKEY = 0
	TUPLEVAL = 1
	TUPLEARITY = 2
)

type PropList erlang.OtpErlangList

func ToMap(list erlang.OtpErlangList) (types.Proplist, error) {
	log.Warning("Preparing to iterate ...")
	tuplesMap := make(types.Proplist)
	for idx, term := range list.Value {
		log.Trace("Iteration ", idx)
		k, v, err := ExtractTuple(term)
		if err != nil {
			log.Error(err)
			return nil, err
		}
		tuplesMap[k] = v
	}
	return tuplesMap, nil
}

func ExtractTuple(term interface{}) (string, interface{}, error) {
	switch t := term.(type) {
	default:
		return "", nil, fmt.Errorf("term %T is not a tuple", t)
	case erlang.OtpErlangTuple:
		tuple := term.(erlang.OtpErlangTuple)
		if len(tuple) != TUPLEARITY {
			return "", nil, fmt.Errorf("tuple of wrong size; expected %d, got %d",
			TUPLEARITY, len(tuple))
		}
		atomKey, ok := tuple[TUPLEKEY].(erlang.OtpErlangAtom)
		if !ok {
			return "", nil, errors.New("unexpected type for directive")
		}
		key := string(atomKey)
		val := tuple[TUPLEVAL]
		return key, val, nil
	}
}