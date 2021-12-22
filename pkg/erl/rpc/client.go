package rpc

type Client struct{}

func New() *Client {
	return new(Client)
}
