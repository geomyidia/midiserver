package datatypes

import (
	"github.com/okeuday/erlang_go/v2/erlang"
	log "github.com/sirupsen/logrus"

	"github.com/ut-proj/midiserver/pkg/erl/packets"
)

func FromPacket(pkt *packets.Packet) (interface{}, error) {
	bytes, err := pkt.Bytes()
	if err != nil {
		log.Error(err)
		return nil, err
	}
	return FromBytes(bytes)
}

func FromBytes(data []byte) (interface{}, error) {
	term, err := erlang.BinaryToTerm(data)
	if err != nil {
		log.Error(err)
		return nil, err
	}
	return FromTerm(term)
}

func FromTerm(term interface{}) (interface{}, error) {
	log.Tracef("Processing term: %+v", term)
	switch t := term.(type) {
	case erlang.OtpErlangAtom:
		log.Debug("Processing atom ...")
		eAtom, ok := term.(erlang.OtpErlangAtom)
		if !ok {
			log.Error(ErrCastingAtom)
			return nil, ErrCastingAtom
		}
		return NewAtom(string(eAtom)), nil
	case erlang.OtpErlangAtomCacheRef:
		log.Debug("Processing atom cache ref ...")
		return nil, ErrNotImplemented
	case erlang.OtpErlangAtomUTF8:
		log.Debug("Processing UTF-8 atom ...")
		return nil, ErrNotImplemented
	case erlang.OtpErlangBinary:
		log.Debug("Processing binary ...")
		eBin, ok := term.(erlang.OtpErlangBinary)
		if !ok {
			log.Error(ErrCastingBinary)
			return nil, ErrCastingBinary
		}
		return NewBinary(eBin.Value), nil
	case erlang.OtpErlangFunction:
		log.Debug("Processing function ...")
		return nil, ErrNotImplemented
	case erlang.OtpErlangList:
		log.Debug("Processing list ...")
		terms := make([]interface{}, len(t.Value))
		for i, term := range t.Value {
			termStruct, err := FromTerm(term)
			if err != nil {
				log.Error(err)
				return nil, ErrCastingTuple
			}
			terms[i] = termStruct
		}
		return NewList(terms), nil
	case erlang.OtpErlangMap:
		log.Debug("Processing map ...")
		return nil, ErrNotImplemented
	case erlang.OtpErlangPid:
		log.Debug("Processing pid ...")
		return nil, ErrNotImplemented
	case erlang.OtpErlangPort:
		log.Debug("Processing port ...")
		return nil, ErrNotImplemented
	case erlang.OtpErlangReference:
		log.Debug("Processing reference ...")
		return nil, ErrNotImplemented
	case erlang.OtpErlangTuple:
		log.Debug("Processing tuple ...")
		terms := make([]interface{}, len(t))
		for i, term := range t {
			termStruct, err := FromTerm(term)
			if err != nil {
				log.Error(err)
				return nil, ErrCastingTuple
			}
			terms[i] = termStruct
		}
		return NewTuple(terms), nil
	// Non-OTP types that are already converted by the OTP library
	case string:
		log.Debug("Processing string ...")
		// return NewAtom(t), nil
		return t, nil
	case int:
		log.Debug("Processing int ...")
		// return NewAtom(strconv.Itoa(t)), nil
		return t, nil
	case byte:
		log.Debug("Processing byte ...")
		// return NewAtom(string(t)), nil
		return t, nil
	case []byte:
		log.Debug("Processing bytes ...")
		// return NewAtom(string(t)), nil
		return t, nil
	default:
		log.Error(ErrUnsupportedOTPType)
		log.Tracef("type: %+v (term: %+v", t, term)
		return nil, ErrUnsupportedOTPType
	}
}
