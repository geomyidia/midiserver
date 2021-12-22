package rpc

import (
	"encoding/binary"
	"errors"
	"io"
	"strings"

	log "github.com/sirupsen/logrus"

	"github.com/ut-proj/midiserver/internal/tcp"
	"github.com/ut-proj/midiserver/pkg/types"
)

const (
	PORT_PLEASE2_REQ = 122
	PORT2_RESP       = 119
	PORT2_OK         = 0
	ErrInt           = -1
	nodeDelimit      = "@"
)

var (
	ErrEPMDEmptyReply = errors.New("empty response from epmd")
	ErrEPMDBadReply   = errors.New("unexpected response from epmd")
	ErrEPMDNotFound   = errors.New("node not found")
)

type Client struct {
	epmdHost   string
	epmdPort   int
	remoteNode string
	remotePort int
	tcpClient  *tcp.Client
}

func New(flags *types.Flags) (*Client, error) {
	client := &Client{
		epmdHost:   flags.EPMDHost,
		epmdPort:   flags.EPMDPort,
		remoteNode: flags.RemoteNode,
	}
	port, err := NodePort(client.epmdHost, client.remotePort, client.remoteNode)
	if err != nil {
		return nil, err
	}
	tcpClient, err := tcp.NewClient(client.remoteNode, port)
	if err != nil {
		return nil, err
	}
	client.remotePort = port
	client.tcpClient = tcpClient
	return client, nil
}

func (c *Client) Close() error {
	return c.tcpClient.Close()
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
	shortName := strings.Split(remoteNode, nodeDelimit)[0]
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
