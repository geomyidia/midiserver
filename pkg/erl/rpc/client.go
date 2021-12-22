package rpc

import (
	log "github.com/sirupsen/logrus"

	"github.com/ut-proj/midiserver/internal/tcp"
	"github.com/ut-proj/midiserver/pkg/types"
)

const (
	intError = -1
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
	port, err := NodePort(client.epmdHost, client.remotePort)
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

func NodePort(epmdHost string, remotePort int) (int, error) {
	epmdClient, err := tcp.NewClient(epmdHost, remotePort)
	if err != nil {
		return intError, err
	}
	payload := ""
	err = epmdClient.WriteStr(payload)
	if err != nil {
		return intError, err
	}
	reply, err := epmdClient.Read()
	if err != nil {
		return intError, err
	}
	parsed, err := parseEPMDReply(reply)
	log.Tracef("parsed data: %v", parsed)
	if err != nil {
		return intError, err
	}
	epmdClient.Close()
	return 0, nil
}

func parseEPMDReply(data []byte) ([]byte, error) {
	return data, nil
}
