package ipc

import (
	"encoding/binary"
	"errors"
	"hash"
	"io"
	"os"
	"os/exec"

	"github.com/davecgh/go-spew/spew"
	"github.com/minio/highwayhash"
)

var quitMessage = ^uint32(0)
var QuitCommand = []byte{255, 255, 255, 255}
var hwhKey = make([]byte, 32)

type Conn struct {
	io.Writer
	io.Reader
	Hash hash.Hash64
	Name string
}

func NewConn(name string, in io.Reader, out io.Writer) (c Conn, err error) {
	c.Name = name
	c.Writer = out
	c.Reader = in
	c.Hash, err = highwayhash.New64(hwhKey)
	if err != nil {
		panic(err)
	}
	return
}

type Controller struct {
	*exec.Cmd
	Conn
}

type Worker struct {
	Conn
}

func NewWorker() (w *Worker, err error) {
	nC, err := NewConn("worker", os.Stdin, os.Stdout)
	if err != nil {
		panic(err)
	}
	w = &Worker{nC}
	return
}

// A controller runs a child process and attaches to its stdin/out
func NewController(args []string) (c *Controller, err error) {
	// if runtime.GOOS == "windows" {
	// 	args = append([]string{"cmd.exe", "/C", "start"}, args...)
	// }
	// args = apputil.PrependForWindows(args)
	c = &Controller{
		Cmd: exec.Command(args[0], args[1:]...),
	}
	var w io.WriteCloser
	if w, err = c.StdinPipe(); Check(err) {
		panic(err)
	}
	// child process can print to parent's stderr for debugging
	c.Cmd.Stderr = os.Stderr
	c.Stderr = os.Stderr
	var r io.ReadCloser
	if r, err = c.StdoutPipe(); Check(err) {
		panic(err)
	}
	if c.Conn, err = NewConn("controller", r, w); Check(err) {
		return
	}
	return
}

// Write sends a message over the IPC pipe containing 32 bit length, 64 bit highway hash and the payload
func (c *Conn) Write(p []byte) (n int, err error) {
	Error("write", spew.Sdump(p))
	pLen := len(p)
	n, err = c.Hash.Write(p)
	if err != nil {
		return
	}
	sum := c.Hash.Sum(nil)
	Error("write", spew.Sdump(sum))
	c.Hash.Reset()
	b := make([]byte, 4)
	binary.LittleEndian.PutUint32(b, uint32(pLen))
	out := append(append(b, sum...), p...)
	Error("write", spew.Sdump(out))
	return c.Writer.Write(out)
}

// Read scans the input for a new message.
//
// First 4 bytes are the payload size in uint32 littleendian,
//
// Second 8 bytes are a highwayhash64 hash of the payload
//
// In the return the decoded length from the header of the incoming bytes which matches the hash, expected message
// length, or errors are returned. the input byte slice must be truncated to the given n or it isn't guaranteed to be
// correct data
func (c *Conn) Read(p []byte) (n int, err error) {
	Error("read buffer length", len(p))
	r := c.Reader.Read
	u64 := binary.LittleEndian.Uint64
	u32 := binary.LittleEndian.Uint32
	n, err = r(p[:4])
	if err != nil || n != 4 {
		Error(err)
	}
	Error("read", p[:4])
	bLen := u32(p[:4])
	Error("read", bLen)
	if bLen == quitMessage {
		return 0, errors.New("quit")
	}
	mHash := p[:8]
	n, err = r(mHash)
	if err != nil || n != 8 {
		return
	}
	Error("read", p[:8])
	hash64 := u64(mHash)
	Error("read", hash64)
	n, err = r(p[:bLen])
	if err != nil {
		return
	}
	if n != int(bLen) {
		return 0, errors.New("short message")
	}
	n, err = c.Hash.Write(p[:bLen])
	if err != nil {
		return
	}
	h := c.Hash.Sum64()
	c.Hash.Reset()
	if h != hash64 {
		return 0, errors.New("corrupted message")
	}
	return
}

// Close signals the worker process to shut down and closes its tty connection
func (c *Conn) Close() (err error) {
	// ciao!
	_, err = c.Write(QuitCommand)
	if err != nil {
		return err
	}
	// hang up
	return
}
