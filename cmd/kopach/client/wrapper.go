package client

import (
	"errors"
	"github.com/p9c/pod/cmd/kopach/control/templates"
	"io"
	"net/rpc"
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
func (c *Client) NewJob(templates *templates.Message) (err error) {
	Trace("sending new templates")
	// Debugs(templates)
	if templates == nil {
		err = errors.New("templates is nil")
		Error(err)
		return
	}
	var reply bool
	if err = c.Call("Worker.NewJob", templates, &reply); Check(err){
		return
	}
	if reply != true {
		err = errors.New("new templates command not acknowledged")
	}
	return
}

// Pause tells the worker to stop working, this is for when the controlling node
// is not current
func (c *Client) Pause() (err error) {
	// Debug("sending pause")
	var reply bool
	err = c.Call("Worker.Pause", 1, &reply)
	if err != nil {
		Error(err)
		return
	}
	if reply != true {
		err = errors.New("pause command not acknowledged")
	}
	return
}

// Stop the workers
func (c *Client) Stop() (err error) {
	Debug("stop working (exit)")
	var reply bool
	err = c.Call("Worker.Stop", 1, &reply)
	if err != nil {
		Error(err)
		return
	}
	if reply != true {
		err = errors.New("stop command not acknowledged")
	}
	return
}

// SendPass sends the multicast PSK to the workers so they can dispatch their
// solutions
func (c *Client) SendPass(pass string) (err error) {
	Debug("sending dispatch password")
	var reply bool
	err = c.Call("Worker.SendPass", pass, &reply)
	if err != nil {
		Error(err)
		return
	}
	if reply != true {
		err = errors.New("send pass command not acknowledged")
	}
	return
}
