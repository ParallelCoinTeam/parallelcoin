package ipclog

import (
	"errors"
	"github.com/p9c/pod/pkg/kopachctrl/job"
	"io"
	"net/rpc"
)

type Client struct {
	*rpc.Client
}

// New creates a new client logi's ipcLogger
// Note that any kind of connection can be used here, other than the StdConn
func New(conn io.ReadWriteCloser) *Client {
	return &Client{rpc.NewClient(conn)}
}

// The following are all blocking calls as they are all triggers rather than
// queries and should return immediately the message is received.
// If deadlines are needed, set them on the connection,
// for StdConn this shouldn't be required as usually if the server is running
// worker will be too, a deadline would be needed for a network connection,
// or alternatively as with the Controller just spew messages over UDP

// Run the delivery of log entries
func (c *Client) Run(job *job.Container) (err error) {
	L.Debug("starting logger feed")
	var reply bool
	err = c.Call("IPCLogger.Run", job, &reply)
	if err != nil {
		L.Error(err)
		return
	}
	if reply != true {
		err = errors.New("start signal not acknowledged")
	}
	return
}

// Pause the delivery of log entries
func (c *Client) Pause(job *job.Container) (err error) {
	L.Debug("stopping logger feed")
	var reply bool
	err = c.Call("IPCLogger.Pause", job, &reply)
	if err != nil {
		L.Error(err)
		return
	}
	if reply == true {
		err = errors.New("stop signal not acknowledged")
	}
	return
}


