package rpc

import (
	"errors"
	"fmt"
	"time"

	"github.com/ergo-services/ergo"
	"github.com/ergo-services/ergo/etf"
	"github.com/ergo-services/ergo/gen"
	"github.com/ergo-services/ergo/node"
	log "github.com/sirupsen/logrus"

	"github.com/ut-proj/midiserver/pkg/erl"
	"github.com/ut-proj/midiserver/pkg/types"
)

const (
	channelBuffers = 2
	clockModule    = "undermidi.extclock"
)

var (
	pingTimeout = time.Second * time.Duration(2)
	ErrTimeout  = errors.New("result timeout")
)

type RPC struct {
	module   string
	function string
	args     []etf.Term
}

type Client struct {
	gen.Server
	node          node.Node
	remoteModule  string
	remoteNode    string
	remoteProcess gen.Process
	result        chan interface{}
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
	remoteProcess, err := remoteNode.Spawn(erl.ShortNodename, gen.ProcessOptions{}, client)
	if err != nil {
		return nil, err
	}
	client.remoteProcess = remoteProcess
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
	val, err := c.send(module, "ping")
	if err != nil {
		return "", err
	}
	return val.(string), nil
}

func (c *Client) ClockTick() (string, error) {
	val, err := c.send(clockModule, "tick")
	if err != nil {
		return "", err
	}
	return val.(string), nil
}

// Private methods

func (c *Client) send(module, function string, args ...etf.Term) (interface{}, error) {
	var err error
	if len(args) == 0 {
		err = c.remoteProcess.Send(c.remoteProcess.Self(), RPC{module: module, function: function})
	} else {
		err = c.remoteProcess.Send(
			c.remoteProcess.Self(), RPC{module: module, function: function, args: args})
	}
	if err != nil {
		return nil, err
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
