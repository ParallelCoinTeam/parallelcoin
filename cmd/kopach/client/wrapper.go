package client

import (
	"errors"
	"github.com/p9c/pod/pkg/chain/wire"
	"github.com/p9c/pod/pkg/log"
	"io"
	"net/rpc"
)

type Client struct {
	*rpc.Client
}

// New creates a new client for a kopach_worker.
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

// NewJob is a delivery of a new job for the worker, this starts a miner
func (w *Client) NewJob(blk *wire.MsgBlock) (err error) {
	log.DEBUG("sending new job")
	var reply bool
	err = w.Call("Worker.NewJob", blk, &reply)
	if err != nil {
		log.ERROR(err)
		return
	}
	if reply != true {
		err = errors.New("new job command not acknowledged")
	}
	return
}

func (w *Client) Pause() (err error) {
	log.DEBUG("sending new job")
	var reply bool
	err = w.Call("Worker.Pause", 1, &reply)
	if err != nil {
		log.ERROR(err)
		return
	}
	if reply != true {
		err = errors.New("pause command not acknowledged")
	}
	return
}

func (w *Client) Stop() (err error) {
	log.DEBUG("sending new job")
	var reply bool
	err = w.Call("Worker.Stop", 1, &reply)
	if err != nil {
		log.ERROR(err)
		return
	}
	if reply != true {
		err = errors.New("stop command not acknowledged")
	}
	return
}
