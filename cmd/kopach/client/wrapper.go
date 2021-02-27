package client

import (
	"errors"
	"io"
	"net/rpc"

	"github.com/p9c/pod/cmd/kopach/control/job"
)

type Client struct {
	*rpc.Client
}

// New creates a new client for a kopach_worker. Note that any kind of connection can be used here, other than the
// StdConn
func New(conn io.ReadWriteCloser) *Client {
	return &Client{rpc.NewClient(conn)}
}

// The following are all blocking calls as they are all triggers rather than queries and should return immediately the
// message is received. If deadlines are needed, set them on the connection, for StdConn this shouldn't be required as
// usually if the server is running worker will be too, a deadline would be needed for a network connection, or
// alternatively as with the Controller just spew messages over UDP

// NewJob is a delivery of a new job for the worker, this starts a miner
func (c *Client) NewJob(job *job.Job) (err error) {
	Trace("sending new job")
	// Debugs(job)
	if job == nil {
		err = errors.New("job is nil")
		Error(err)
		return
	}
	var reply bool
	if err = c.Call("Worker.NewJob", job, &reply); Check(err){
		return
	}
	if reply != true {
		err = errors.New("new job command not acknowledged")
	}
	return
}

// Pause tells the worker to stop working, this is for when the controlling node is not current
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
