package rpc

import (
	"encoding/binary"
	"errors"
	"fmt"
	"io"
	"strings"
	"time"

	"github.com/ergo-services/ergo"
	"github.com/ergo-services/ergo/etf"
	"github.com/ergo-services/ergo/gen"
	"github.com/ergo-services/ergo/node"
	log "github.com/sirupsen/logrus"

	"github.com/ut-proj/midiserver/internal/tcp"
	"github.com/ut-proj/midiserver/pkg/erl"
	"github.com/ut-proj/midiserver/pkg/types"
)

const (
	PORT_PLEASE2_REQ = 122
	PORT2_RESP       = 119
	PORT2_OK         = 0
	ErrInt           = -1
	channelBuffers   = 2
)

var (
	pingTimeout = time.Second * time.Duration(2)
	// Errors
	ErrEPMDEmptyReply = errors.New("empty response from epmd")
	ErrEPMDBadReply   = errors.New("unexpected response from epmd")
	ErrEPMDNotFound   = errors.New("node not found")
	ErrBadNodeName    = errors.New("bad node name")
	ErrTimeout        = errors.New("result timeout")
)

type RPC struct {
	module   string
	function string
	args     []etf.Term
}

type Client struct {
	gen.Server
	node         node.Node
	remoteNode   string
	remoteModule string
	result       chan interface{}
}

func New(flags *types.Flags) (*Client, error) {
	cookie, err := erl.ReadCookie()
	if err != nil {
		return nil, err
	}
	remoteNode, err := ergo.StartNode(erl.LongNodename, cookie, node.Options{})
	if err != nil {
		return nil, err
	}
	client := &Client{
		node:         remoteNode,
		result:       make(chan interface{}, channelBuffers),
		remoteNode:   flags.RemoteNode,
		remoteModule: flags.RemoteModule,
	}

	return client, nil
}

func (c *Client) HandleInfo(process *gen.ServerProcess, message etf.Term) gen.ServerStatus {
	switch msg := message.(type) {
	case RPC:
		var value etf.Term
		var err error
		log.Debug("Got RPC message ...")
		if len(msg.args) == 0 {
			value, err = process.CallRPC(c.remoteNode, msg.module, msg.function)
		} else {
			value, err = process.CallRPC(c.remoteNode, msg.module, msg.function, msg.args...)
		}
		if err != nil {
			c.result <- err
		} else {
			c.result <- value
		}
		return gen.ServerStatusOK
	}
	return gen.ServerStatusStop
}

func (c *Client) Ping(module string) (string, error) {
	process, err := c.node.Spawn(erl.ShortNodename, gen.ProcessOptions{}, c)
	if err != nil {
		return "", err
	}
	err = process.Send(process.Self(), RPC{module: module, function: "ping"})
	if err != nil {
		return "", err
	}
	result, err := c.awaitResult()
	if err != nil {
		return "", err
	}
	return result.(string), nil
}

func (c *Client) awaitResult() (interface{}, error) {
	select {
	case value := <-c.result:
		log.Debugf("got value: %v", value)
		switch val := value.(type) {
		case etf.Atom:
			return string(val), nil
		case error:
			return nil, val
		default:
			return nil, fmt.Errorf("unexpected value type: %v", val)
		}

	case <-time.After(pingTimeout):
		return "", ErrTimeout
	}
}

// Utility functions

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
	parsed, err := parseEPMDReply(reply)
	log.Tracef("parsed data: %v", parsed)
	if err != nil {
		return ErrInt, err
	}
	return parsed, nil
}

// Taken from https://github.com/ergo-services/ergo/blob/master/node/epmd.go
func makePortPlease2Req(remoteNode string) []byte {
	shortName := strings.Split(remoteNode, erl.NodeDelimit)[0]
	reqLen := uint16(2 + len(shortName) + 1)
	req := make([]byte, reqLen)
	binary.BigEndian.PutUint16(req[0:2], uint16(len(req)-2))
	req[2] = byte(PORT_PLEASE2_REQ)
	copy(req[3:reqLen], shortName)
	return req
}

func parseEPMDReply(data []byte) (int, error) {
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

func bytesToInt(b []byte) uint16 {
	return binary.BigEndian.Uint16(b)
}
