package tcp

import (
	"fmt"
	"io"
	"net"

	log "github.com/sirupsen/logrus"
)

const (
	defaultProto = "tcp"
)

type Client struct {
	conn *net.TCPConn
}

func NewClient(host string, port int) (*Client, error) {
	hostPort := fmt.Sprintf("%s:%d", host, port)
	tcpAddr, err := net.ResolveTCPAddr(defaultProto, hostPort)
	if err != nil {
		return nil, err
	}

	conn, err := net.DialTCP(defaultProto, nil, tcpAddr)
	if err != nil {
		return nil, err
	}
	return &Client{conn: conn}, nil
}

func (c *Client) WriteStr(data string) error {
	return c.Write([]byte(data))
}

func (c *Client) Write(data []byte) error {
	log.Tracef("writing data to TCP connection: %v (%s)", data, string(data))
	_, err := c.conn.Write(data)
	return err
}

func (c *Client) Read() ([]byte, error) {
	reply := make([]byte, 1024)
	_, err := c.conn.Read(reply)
	if err == io.EOF {
		return reply, nil
	}
	if err != nil {
		return nil, err
	}
	return reply, nil
}

func (c *Client) Close() error {
	return c.conn.Close()
}
