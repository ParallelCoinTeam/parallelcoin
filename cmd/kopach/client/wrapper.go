package client

import (
	"errors"
	"io"
	"net/rpc"

	"github.com/p9c/matrjoska/pkg/chainrpc/templates"
)

type Client struct {
	*rpc.Client
}

// New creates a new client for a kopach_worker. Note that any kind of connection can be used here, other than the
// StdConn
func New(conn io.ReadWriteCloser) *Client {
	return &Client{rpc.NewClient(conn)}
}

// NewJob is a delivery of a new job for the worker, this starts a miner
// note that since this implements net/rpc by default this is gob encoded
func (c *Client) NewJob(templates *templates.Message) (e error) {
	// T.Ln("sending new templates")
	// D.S(templates)
	if templates == nil {
		e = errors.New("templates is nil")
		return
	}
	var reply bool
	if e = c.Call("Worker.NewJob", templates, &reply); E.Chk(e) {
		return
	}
	if reply != true {
		e = errors.New("new templates command not acknowledged")
	}
	D.Ln("new job delivered to workers")
	return
}

// Pause tells the worker to stop working, this is for when the controlling node
// is not current
func (c *Client) Pause() (e error) {
	// D.Ln("sending pause")
	var reply bool
	e = c.Call("Worker.Pause", 1, &reply)
	if e != nil {
		return
	}
	if reply != true {
		e = errors.New("pause command not acknowledged")
	}
	return
}

// Stop the workers
func (c *Client) Stop() (e error) {
	D.Ln("stop working (exit)")
	var reply bool
	e = c.Call("Worker.Stop", 1, &reply)
	if e != nil {
		return
	}
	if reply != true {
		e = errors.New("stop command not acknowledged")
	}
	return
}

// SendPass sends the multicast PSK to the workers so they can dispatch their
// solutions
func (c *Client) SendPass(pass []byte) (e error) {
	D.Ln("sending dispatch password")
	var reply bool
	e = c.Call("Worker.SendPass", pass, &reply)
	if e != nil {
		return
	}
	if reply != true {
		e = errors.New("send pass command not acknowledged")
	}
	return
}
