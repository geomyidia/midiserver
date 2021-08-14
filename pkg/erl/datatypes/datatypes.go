package datatypes

import (
	"errors"
	"fmt"

	erlang "github.com/okeuday/erlang_go/v2/erlang"
	log "github.com/sirupsen/logrus"

	"github.com/geomyidia/midiserver/pkg/types"
)

const (
	TUPLEKEY   = 0
	TUPLEVAL   = 1
	TUPLEARITY = 2
)

type PropList erlang.OtpErlangList

func PropListToMap(list erlang.OtpErlangList) (types.PropList, error) {
	log.Debug("converting proplist ...")
	tuplesMap := make(types.PropList)
	for idx, term := range list.Value {
		k, v, err := Tuple(term)
		if err != nil {
			log.Error(err)
			return nil, err
		}
		tuplesMap[k] = v
		log.Tracef("  element %d done (%+v).", idx, v)
	}
	return tuplesMap, nil
}

func Tuple(term interface{}) (string, interface{}, error) {
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

func TupleHasKey(term interface{}, soughtKey string) bool {
	key, _, err := Tuple(term)
	if err != nil {
		log.Error(err)
		return false
	}
	return key == soughtKey
}
