package hello

import (
	"io"
	"net/rpc"
)

type Client struct {
	*rpc.Client
}

func NewClient(conn io.ReadWriteCloser) *Client {
	return &Client{rpc.NewClient(conn)}

}

func (h *Client) Say(name string) (reply string) {
	e := h.Call("Hello.Say", "worker", &reply)
	if e != nil  {
				return "error: " + e.Error()
	}
	return
}

func (h *Client) Bye() (reply string) {
	e := h.Call("Hello.Bye", 1, &reply)
	if e != nil  {
				return "error: " + e.Error()
	}
	return
}
