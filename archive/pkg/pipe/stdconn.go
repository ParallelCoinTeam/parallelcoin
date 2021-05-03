package pipe

import (
	"fmt"
	"io"
	"net"
	"runtime"
	"time"
	
	"github.com/p9c/qu"
)

type StdConn struct {
	io.ReadCloser
	io.WriteCloser
	Quit qu.C
}

func New(in io.ReadCloser, out io.WriteCloser, quit qu.C) (s *StdConn) {
	s = &StdConn{in, out, quit}
	_, file, line, _ := runtime.Caller(1)
	o := fmt.Sprintf("%s:%d", file, line)
	T.Ln("new StdConn at", o)
	return
}

func (s *StdConn) Read(b []byte) (n int, e error) {
	return s.ReadCloser.Read(b)
}

func (s *StdConn) Write(b []byte) (n int, e error) {
	return s.WriteCloser.Write(b)
}

func (s *StdConn) Close() (e error) {
	s.Quit.Q()
	return
}

func (s *StdConn) LocalAddr() (addr net.Addr) {
	// this is a no-op as it is not relevant to the type of connection
	return
}

func (s *StdConn) RemoteAddr() (addr net.Addr) {
	// this is a no-op as it is not relevant to the type of connection
	return
}

func (s *StdConn) SetDeadline(t time.Time) (e error) {
	// this is a no-op as it is not relevant to the type of connection
	return
}

func (s *StdConn) SetReadDeadline(t time.Time) (e error) {
	// this is a no-op as it is not relevant to the type of connection
	return
}

func (s *StdConn) SetWriteDeadline(t time.Time) (e error) {
	// this is a no-op as it is not relevant to the type of connection
	return
}
