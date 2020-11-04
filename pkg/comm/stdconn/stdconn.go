package stdconn

import (
	"io"
	"net"
	"time"
)

type StdConn struct {
	io.ReadCloser
	io.WriteCloser
	Quit chan struct{}
}

func New(in io.ReadCloser, out io.WriteCloser, quit chan struct{}) (s StdConn) {
	s = StdConn{in, out, quit}
	go func() {
	out:
		for {
			select {
			case <-quit:
				if err := s.ReadCloser.Close(); Check(err) {
				}
				if err := s.WriteCloser.Close(); Check(err) {
				}
				break out
			}
		}
	}()
	return
}

func (s StdConn) Read(b []byte) (n int, err error) {
	return s.ReadCloser.Read(b)
}

func (s StdConn) Write(b []byte) (n int, err error) {
	return s.WriteCloser.Write(b)
}

func (s StdConn) Close() (err error) {
	close(s.Quit)
	return
}

func (s StdConn) LocalAddr() (addr net.Addr) {
	// this is a no-op as it is not relevant to the type of connection
	return
}

func (s StdConn) RemoteAddr() (addr net.Addr) {
	// this is a no-op as it is not relevant to the type of connection
	return
}

func (s StdConn) SetDeadline(t time.Time) (err error) {
	// this is a no-op as it is not relevant to the type of connection
	return
}

func (s StdConn) SetReadDeadline(t time.Time) (err error) {
	// this is a no-op as it is not relevant to the type of connection
	return
}

func (s StdConn) SetWriteDeadline(t time.Time) (err error) {
	// this is a no-op as it is not relevant to the type of connection
	return
}
