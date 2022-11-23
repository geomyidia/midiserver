package datatypes

import (
	"github.com/okeuday/erlang_go/v2/erlang"
	log "github.com/sirupsen/logrus"
)

func TupleToSlice(tuple erlang.OtpErlangTuple) ([]string, error) {
	keyTerm := tuple[tupleKey]
	key, ok := keyTerm.(erlang.OtpErlangAtom)
	if !ok {
		log.Error(ErrCastingAtom)
		return nil, ErrCastingAtom
	}
	valTerm := tuple[tupleVal]
	val, ok := valTerm.(erlang.OtpErlangAtom)
	if !ok {
		log.Error(ErrCastingAtom)
		return nil, ErrCastingAtom
	}
	return []string{string(key), string(val)}, nil
}

func TupleListToMap(list erlang.OtpErlangList) (map[string]string, error) {
	m := make(map[string]string)
	for _, t := range list.Value {
		tpl, ok := t.(erlang.OtpErlangTuple)
		if !ok {
			log.Error(ErrCastingTuple)
			return nil, ErrCastingTuple
		}
		slice, err := TupleToSlice(tpl)
		if err != nil {
			return nil, err
		}
		m[slice[tupleKey]] = slice[tupleVal]
	}
	return m, nil
}

func MapStrsToInterfaces(oldMap map[string]string) map[string]interface{} {
	newMap := make(map[string]interface{})
	for k, v := range oldMap {
		newMap[k] = v
	}
	return newMap
}
