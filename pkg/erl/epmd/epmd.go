package epmd

import (
	"bytes"
	"encoding/binary"
	"errors"
	"io"
	"strings"

	log "github.com/sirupsen/logrus"

	"github.com/ut-proj/midiserver/internal/tcp"
	"github.com/ut-proj/midiserver/pkg/erl"
)

const (
	EPMD_NAMES_REQ   = 110
	PORT_PLEASE2_REQ = 122
	PORT2_RESP       = 119
	PORT2_OK         = 0
	ErrInt           = -1
	ErrStr           = "ERR"
)

var (
	EPMD_NAMES_RESP = []byte{0, 0, 17, 17}
	ErrStrSlice     = []string{ErrStr}
	// Errors
	ErrEPMDEmptyReply = errors.New("empty response from epmd")
	ErrEPMDBadReply   = errors.New("unexpected response from epmd")
	ErrEPMDNotFound   = errors.New("node not found")
	ErrBadNodeName    = errors.New("bad node name")
)

func ListNodes(epmdHost string, epmdPort int) ([]string, error) {
	epmdClient, err := tcp.NewClient(epmdHost, epmdPort)
	if err != nil {
		return ErrStrSlice, err
	}
	defer epmdClient.Close()
	err = epmdClient.Write(makeEpmdNamesReq())
	if err != nil {
		return ErrStrSlice, err
	}
	reply, err := epmdClient.Read()
	if err != nil && err != io.EOF {
		return ErrStrSlice, err
	}
	if len(reply) == 0 {
		return ErrStrSlice, ErrEPMDEmptyReply
	}
	log.Tracef("raw reply: %v", reply)
	log.Tracef("raw reply (string): %s", string(reply))
	parsed, err := parseNodesReply(reply)
	log.Tracef("parsed data: %q", parsed)
	if err != nil {
		return ErrStrSlice, err
	}
	return parsed, nil
}

func NodePort(epmdHost string, epmdPort int, remoteNode string) (int, error) {
	epmdClient, err := tcp.NewClient(epmdHost, epmdPort)
	if err != nil {
		return ErrInt, err
	}
	defer epmdClient.Close()
	err = epmdClient.Write(makePortPlease2Req(remoteNode))
	if err != nil {
		return ErrInt, err
	}
	reply, err := epmdClient.Read()
	if err != nil && err != io.EOF {
		return ErrInt, err
	}
	if len(reply) == 0 {
		return ErrInt, ErrEPMDEmptyReply
	}
	log.Tracef("raw reply: %v", reply)
	parsed, err := parsePortReply(reply)
	log.Tracef("parsed data: %v", parsed)
	if err != nil {
		return ErrInt, err
	}
	return parsed, nil
}

func makeEpmdNamesReq() []byte {
	return makeCodeReq(EPMD_NAMES_REQ)
}

func makePortPlease2Req(remoteNode string) []byte {
	data := shortName(remoteNode)
	return makeReq(data, PORT_PLEASE2_REQ)
}

func parseNodesReply(data []byte) ([]string, error) {
	respCode := data[0:4]
	log.Tracef("Resp code: %v", respCode)
	if !bytes.Equal(respCode, EPMD_NAMES_RESP) {
		return ErrStrSlice, ErrEPMDBadReply
	}
	respData := string(data[4 : len(data)-1])
	log.Tracef("Resp data: %v", respData)
	return strings.Split(respData, "\n"), nil
}

func parsePortReply(data []byte) (int, error) {
	distCode := data[0]
	if distCode != PORT2_RESP {
		return ErrInt, ErrEPMDBadReply
	}
	resultCode := data[1]
	if resultCode != PORT2_OK {
		return ErrInt, ErrEPMDNotFound
	}
	port := bytesToInt(data[2:4])
	return int(port), nil
}

func makeCodeReq(reqCode int) []byte {
	return makeReq("", reqCode)
}

// Taken from https://github.com/ergo-services/ergo/blob/master/node/epmd.go
func makeReq(data string, reqCode int) []byte {
	dataLen := len(data)
	reqLen := uint16(2 + dataLen + 1)
	req := make([]byte, reqLen)
	binary.BigEndian.PutUint16(req[0:2], uint16(len(req)-2))
	req[2] = byte(reqCode)
	if dataLen > 0 {
		copy(req[3:reqLen], data)
	}
	return req
}

func shortName(remoteNode string) string {
	return strings.Split(remoteNode, erl.NodeDelimit)[0]
}

func bytesToInt(b []byte) uint16 {
	return binary.BigEndian.Uint16(b)
}
