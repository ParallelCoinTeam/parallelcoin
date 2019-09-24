package client

import (
	"net"

	"github.com/p9c/kcp"
)

func newDirectKCPConn(c *Client, network, address string) (net.Conn, error) {
	return kcp.DialWithOptions(address, c.option.Block.(kcp.BlockCrypt), 10, 3)
}
