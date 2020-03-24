package serve

import (
	"errors"
	"github.com/p9c/logi"
	"io"
	"net/rpc"
)

type Serve struct {
	io.ReadWriteCloser
	*rpc.Client
}

// New creates a new provider of logi.Entry records from the logger
func New(conn io.ReadWriteCloser) *Serve {
	s := &Serve{
		ReadWriteCloser: conn,
	}
	s.Client = rpc.NewClient(s)
	return s
}

// Log sends a new log message
func (c *Serve) Log(ent logi.Entry) (err error) {
	L.Debug("starting logger feed")
	var reply bool
	err = c.Call("API.Log", ent, &reply)
	if err != nil {
		L.Error(err)
		return
	}
	if reply != true {
		err = errors.New("message not acknowledged")
	}
	return
}

// API is the public API for serving logs to a consume.API
type API struct {
	Serve  *Serve
	RunC   chan struct{}
	PauseC chan struct{}
	EntryC chan logi.Entry
	Quit   chan struct{}
}

// NewAPI creates a new API server
func NewAPI(conn io.ReadWriteCloser) *API {
	a := &API{
		Serve:  New(conn),
		RunC:   make(chan struct{}),
		PauseC: make(chan struct{}),
	}
	if err := rpc.Register(a); L.Check(err) {
	}
	go func() {
		L.Debug("starting up Serve IPC")
		rpc.ServeConn(conn)
		L.Debug("stopping Serve IPC")
		if err := conn.Close(); L.Check(err) {
		}
		L.Debug("finished Serve")
	}()
	return a
}

// Run starts serving Log
func (a *API) Run(cmd bool, reply *bool) (err error) {
	L.Debug("serving logs")
	a.RunC <- struct{}{}
	r := true
	reply = &r
	return
}

// Pause pauses serving Log
func (a *API) Pause(cmd bool, reply *bool) (err error) {
	L.Debug("pausing log service")
	a.PauseC <- struct{}{}
	r := false
	reply = &r
	return
}

// Run is the loop that responds to the run and pause messages and receives
// entries to send
func Run(a *API) (err error) {
out:
	for {
		select {
		case <-a.Quit:
			break out
		case <-a.EntryC:
		case <-a.RunC:
		running:
			for {
				select {
				case <-a.Quit:
					break out
				case e := <-a.EntryC:
					if err := a.Serve.Log(e); L.Check(err) {
					}
				case <-a.RunC:
				case <-a.PauseC:
					break running
				}
			}
		case <-a.PauseC:
		pausing:
			for {
				select {
				case <-a.Quit:
					break out
				case <-a.EntryC:
				case <-a.RunC:
					break pausing
				case <-a.PauseC:
				}
			}
		}
	}
	return
}
