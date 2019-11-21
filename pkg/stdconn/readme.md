# StdConn

This is a simple library that runs a client executable (in the example, it
uses go run), creates a net.Conn with the stdin/stdout on both sides so the
controller process can send IPC messages

## How to use

`stdconn.go` and `worker/worker.go` are the reusable/convenience parts
. `stdconn.go` is the net.Conn implementation and `worker.go` is mainly just
 a helper to start up an arbitrary executable with arbitrary arguments, it
  also has helpers to send stop and kill signals.
 
`hello/implementation.go` contains an example RPC implementation. Note that the
most basic implementation should have a function like `Bye()` to stop the
client, as closing its connection when using this IPC mechanism does not stop 
the child process which will block unless it is explicitly killed, and using
os.Signal requires adding this handling code. So they are there for defcon IV
type situations mainly and if using signals is preferred.
 
`hello/wrapper` contains a set of wrappers that shortens the RPC call syntax.
If asynchronous calling is required, the wrappers would mimic the calling
pattern with a cancel function and function that returns the reply or error
over a channel. Probably this could be automated somehow but I haven't got
the need or time for that right now. 