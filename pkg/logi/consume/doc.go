// Package consume turns on and off log entry messages from a Serve endpoint
//
// For this all you need is a handler function and to call x.Consume.Run() will
// call your handler when a log arrives. To start logging call Run() and to pouse
// it call Pause()
package consume
