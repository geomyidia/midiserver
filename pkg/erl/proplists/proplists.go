package proplists

import (
	erlang "github.com/okeuday/erlang_go/v2/erlang"
	log "github.com/sirupsen/logrus"
)

const (
	TUPLEKEY = 0
	TUPLEVAL = 1
)

type PropList erlang.OtpErlangList

func ToMap(list erlang.OtpErlangList) map[string]interface{} {
	log.Warning("Preparing to iterate ...")
	tuplesMap := make(map[string]interface{})
	for idx, term := range list.Value {
		log.Trace("Iteration ", idx)
		switch t := term.(type) {
		case erlang.OtpErlangTuple:
			k, v := ExtractTuple(term)
			tuplesMap[k] = v
			log.Warning("do a thing with the tuple")
		default:
			log.Warningf("unexpected type %T", t)
		}
	}
	return tuplesMap
}

func ExtractTuple(term interface{}) (string, interface{}) {
	tuple := term.(erlang.OtpErlangTuple)
	key := string(tuple[TUPLEKEY].(erlang.OtpErlangAtom))
	val := tuple[TUPLEVAL]
	return key, val
}
