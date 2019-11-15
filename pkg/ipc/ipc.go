package ipc

import (
	"io"
	"os"
	"os/exec"
)

var QuitCommand = []byte{255, 255, 255, 255}

type Controller struct {
	*exec.Cmd
	In  io.Writer
	Out io.Reader
}

func NewController() (out *Controller, err error) {
	l := len(os.Args)
	args := append(os.Args[:l-1], "worker")
	out = &Controller{
		Cmd: exec.Command(args[0], args[1:]...),
	}
	out.In, err = out.StdinPipe()
	if err != nil {
		panic(err)
	}
	out.Stderr = os.Stdout
	out.Out, err = out.StdoutPipe()
	if err != nil {
		panic(err)
	}
	return
}

func (c *Controller) Write(p []byte)  (n int, err error) {
	return c.In.Write(p)
}

func (c *Controller) Read(p []byte) (n int, err error) {
	return c.Out.Read(p)
}

func (c *Controller) Close() error {
	return c.Close()
}

type Worker struct {
	In  io.Writer
	Out io.Reader
}
