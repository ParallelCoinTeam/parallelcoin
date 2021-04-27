package rpcclient

import (
	js "encoding/json"
	"errors"
	
	"github.com/p9c/matrjoska/pkg/btcjson"
)

// FutureRawResult is a future promise to deliver the result of a RawRequest RPC invocation (or an applicable error).
type FutureRawResult chan *response

// Receive waits for the response promised by the future and returns the raw response, or an error if the request was unsuccessful.
func (r FutureRawResult) Receive() (js.RawMessage, error) {
	return receiveFuture(r)
}

// RawRequestAsync returns an instance of a type that can be used to get the result of a custom RPC request at some
// future time by invoking the Receive function on the returned instance.
//
// See RawRequest for the blocking version and more details.
func (c *Client) RawRequestAsync(method string, params []js.RawMessage) FutureRawResult {
	// Method may not be empty.
	if method == "" {
		return newFutureError(errors.New("no method"))
	}
	// Marshal parameters as "[]" instead of "null" when no parameters are passed.
	if params == nil {
		params = []js.RawMessage{}
	}
	// Create a raw JSON-RPC request using the provided method and netparams and marshal it. This is done rather than
	// using the sendCmd function since that relies on marshalling registered json commands rather than custom commands.
	id := c.NextID()
	rawRequest := &btcjson.Request{
		Jsonrpc: "1.0",
		ID:      id,
		Method:  method,
		Params:  params,
	}
	marshalledJSON, e := js.Marshal(rawRequest)
	if e != nil {
		return newFutureError(e)
	}
	// Generate the request and send it along with a channel to respond on.
	responseChan := make(chan *response, 1)
	jReq := &jsonRequest{
		id:             id,
		method:         method,
		cmd:            nil,
		marshalledJSON: marshalledJSON,
		responseChan:   responseChan,
	}
	c.sendRequest(jReq)
	return responseChan
}

// RawRequest allows the caller to send a raw or custom request to the server.
//
// This method may be used to send and receive requests and responses for requests that are not handled by this client
// package, or to proxy partially unmarshalled requests to another JSON-RPC server if a request cannot be handled
// directly.
func (c *Client) RawRequest(method string, params []js.RawMessage) (js.RawMessage, error) {
	return c.RawRequestAsync(method, params).Receive()
}
