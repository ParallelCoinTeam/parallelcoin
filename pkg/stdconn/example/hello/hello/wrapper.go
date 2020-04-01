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
	err := h.Call("Hello.Say", "worker", &reply)
	if err != nil {
		Error(err)
		return "error: " + err.Error()
	}
	return
}

func (h *Client) Bye() (reply string) {
	err := h.Call("Hello.Bye", 1, &reply)
	if err != nil {
		Error(err)
		return "error: " + err.Error()
	}
	return
}
