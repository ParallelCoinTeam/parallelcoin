package consume

import (
	"errors"
	"github.com/p9c/pod/pkg/logi"
	"io"
	"net/rpc"
)

type Consume struct {
	Conn io.ReadWriteCloser
	*rpc.Client
}

// New creates a new client to listen to logs
// Note that any kind of connection can be used here, other than the StdConn
func New(conn io.ReadWriteCloser) *Consume {
	c := &Consume{
		Conn: conn,
	}
	c.Client = rpc.NewClient(conn)
	return c
}

// The following are all blocking calls as they are all triggers rather than
// queries and should return immediately the message is received.
// If deadlines are needed, set them on the connection,
// for StdConn this shouldn't be required as usually if the server is running
// worker will be too, a deadline would be needed for a network connection,
// or alternatively as with the Controller just spew messages over UDP

// Run the delivery of log entries
func (c *Consume) Run() (err error) {
	L.Debug("starting logger feed")
	var reply bool
	err = c.Call("API.Run", true, &reply)
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
func (c *Consume) Pause() (err error) {
	L.Debug("stopping logger feed")
	var reply bool
	err = c.Call("API.Pause", false, &reply)
	if err != nil {
		L.Error(err)
		return
	}
	if reply == true {
		err = errors.New("pause signal not acknowledged")
	}
	return
}

// API is the public API for listening to a serve.API
type API struct {
	Consume *Consume
	Handler func(ent logi.Entry) error
}

// NewAPI creates a new API server
func NewAPI(conn io.ReadWriteCloser, handler func(ent logi.Entry) error) *API {
	a := &API{
		Consume: New(conn),
		Handler: handler,
	}
	if err := rpc.Register(a); L.Check(err) {
	}
	go func() {
		L.Debug("starting up Consume IPC")
		rpc.ServeConn(conn)
		L.Debug("stopping Consume IPC")
		if err := conn.Close(); L.Check(err) {
		}
		L.Debug("finished Consume")
	}()
	return a
}

// Log receives a log message
func (a *API) Log(ent logi.Entry, reply *bool) (err error) {
	L.Debug("log entry received")
	if a.Handler != nil {
		if err = a.Handler(ent); L.Check(err) {
		}
	}
	r := true
	reply = &r
	return
}
