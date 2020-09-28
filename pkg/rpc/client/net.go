package rpcclient

import (
	js "encoding/json"
	"github.com/p9c/pkg/app/slog"

	"github.com/p9c/pod/pkg/rpc/btcjson"
)

// AddNodeCommand enumerates the available commands that the AddNode function accepts.
type AddNodeCommand string

// Constants used to indicate the command for the AddNode function.
const (
	// ANAdd indicates the specified host should be added as a persistent peer.
	ANAdd AddNodeCommand = "add"
	// ANRemove indicates the specified peer should be removed.
	ANRemove AddNodeCommand = "remove"
	// ANOneTry indicates the specified host should try to connect once, but it should not be made persistent.
	ANOneTry AddNodeCommand = "onetry"
)

// String returns the AddNodeCommand in human-readable form.
func (cmd AddNodeCommand) String() string {
	return string(cmd)
}

// FutureAddNodeResult is a future promise to deliver the result of an AddNodeAsync RPC invocation (or an applicable error).
type FutureAddNodeResult chan *response

// Receive waits for the response promised by the future and returns an error if any occurred when performing the specified command.
func (r FutureAddNodeResult) Receive() (err error) {
	_, err = receiveFuture(r)
	return err
}

// AddNodeAsync returns an instance of a type that can be used to get the result of the RPC at some future time by invoking the Receive function on the returned instance. See AddNode for the blocking version and more details.
func (c *Client) AddNodeAsync(host string, command AddNodeCommand) FutureAddNodeResult {
	cmd := btcjson.NewAddNodeCmd(host, btcjson.AddNodeSubCmd(command))
	return c.sendCmd(cmd)
}

// AddNode attempts to perform the passed command on the passed persistent peer. For example, it can be used to add or a remove a persistent peer, or to do a one time connection to a peer. It may not be used to remove non-persistent peers.
func (c *Client) AddNode(host string, command AddNodeCommand) (err error) {
	return c.AddNodeAsync(host, command).Receive()
}

// FutureNodeResult is a future promise to deliver the result of a NodeAsync RPC invocation (or an applicable error).
type FutureNodeResult chan *response

// Receive waits for the response promised by the future and returns an error if any occurred when performing the specified command.
func (r FutureNodeResult) Receive() (err error) {
	_, err = receiveFuture(r)
	return err
}

// NodeAsync returns an instance of a type that can be used to get the result of the RPC at some future time by invoking the Receive function on the returned instance. See Node for the blocking version and more details.
func (c *Client) NodeAsync(command btcjson.NodeSubCmd, host string,
	connectSubCmd *string) FutureNodeResult {
	cmd := btcjson.NewNodeCmd(command, host, connectSubCmd)
	return c.sendCmd(cmd)
}

// Node attempts to perform the passed node command on the host. For example, it can be used to add or a remove a persistent peer, or to do connect or diconnect a non-persistent one. The connectSubCmd should be set either "perm" or "temp", depending on whether we are targetting a persistent or non-persistent peer. Passing nil will cause the default value to be used, which currently is "temp".
func (c *Client) Node(command btcjson.NodeSubCmd, host string,
	connectSubCmd *string) (err error) {
	return c.NodeAsync(command, host, connectSubCmd).Receive()
}

// FutureGetAddedNodeInfoResult is a future promise to deliver the result of a GetAddedNodeInfoAsync RPC invocation (or an applicable error).
type FutureGetAddedNodeInfoResult chan *response

// Receive waits for the response promised by the future and returns information about manually added (persistent) peers.
func (r FutureGetAddedNodeInfoResult) Receive() (nodeInfo []btcjson.GetAddedNodeInfoResult, err error) {
	var res []byte
	if res, err = receiveFuture(r); slog.Check(err) {
		return
	}
	// Unmarshal as an array of getaddednodeinfo result objects.
	if err = js.Unmarshal(res, &nodeInfo); slog.Check(err) {
		return
	}
	return
}

// GetAddedNodeInfoAsync returns an instance of a type that can be used to get the result of the RPC at some future time by invoking the Receive function on the returned instance. See GetAddedNodeInfo for the blocking version and more details.
func (c *Client) GetAddedNodeInfoAsync(peer string) FutureGetAddedNodeInfoResult {
	cmd := btcjson.NewGetAddedNodeInfoCmd(true, &peer)
	return c.sendCmd(cmd)
}

// GetAddedNodeInfo returns information about manually added (persistent) peers. See GetAddedNodeInfoNoDNS to retrieve only a list of the added (persistent) peers.
func (c *Client) GetAddedNodeInfo(peer string) (res []btcjson.GetAddedNodeInfoResult, err error) {
	return c.GetAddedNodeInfoAsync(peer).Receive()
}

// FutureGetAddedNodeInfoNoDNSResult is a future promise to deliver the result of a GetAddedNodeInfoNoDNSAsync RPC invocation (or an applicable error).
type FutureGetAddedNodeInfoNoDNSResult chan *response

// Receive waits for the response promised by the future and returns a list of manually added (persistent) peers.
func (r FutureGetAddedNodeInfoNoDNSResult) Receive() (nodes []string, err error) {
	var res []byte
	if res, err = receiveFuture(r); slog.Check(err) {
		return
	}
	// Unmarshal result as an array of strings.
	if err = js.Unmarshal(res, &nodes); slog.Check(err) {
		return
	}
	return
}

// GetAddedNodeInfoNoDNSAsync returns an instance of a type that can be used to get the result of the RPC at some future time by invoking the Receive function on the returned instance. See GetAddedNodeInfoNoDNS for the blocking version and more details.
func (c *Client) GetAddedNodeInfoNoDNSAsync(peer string) FutureGetAddedNodeInfoNoDNSResult {
	cmd := btcjson.NewGetAddedNodeInfoCmd(false, &peer)
	return c.sendCmd(cmd)
}

// GetAddedNodeInfoNoDNS returns a list of manually added (persistent) peers. This works by setting the dns flag to false in the underlying RPC. See GetAddedNodeInfo to obtain more information about each added (persistent) peer.
func (c *Client) GetAddedNodeInfoNoDNS(peer string) (res []string, err error) {
	return c.GetAddedNodeInfoNoDNSAsync(peer).Receive()
}

// FutureGetConnectionCountResult is a future promise to deliver the result of a GetConnectionCountAsync RPC invocation (or an applicable error).
type FutureGetConnectionCountResult chan *response

// Receive waits for the response promised by the future and returns the number of active connections to other peers.
func (r FutureGetConnectionCountResult) Receive() (count int64, err error) {
	var res []byte
	if res, err = receiveFuture(r); slog.Check(err) {
		return
	}
	// Unmarshal result as an int64.
	if err = js.Unmarshal(res, &count); slog.Check(err) {
		return
	}
	return
}

// GetConnectionCountAsync returns an instance of a type that can be used to get the result of the RPC at some future time by invoking the Receive function on the returned instance. See GetConnectionCount for the blocking version and more details.
func (c *Client) GetConnectionCountAsync() FutureGetConnectionCountResult {
	cmd := btcjson.NewGetConnectionCountCmd()
	return c.sendCmd(cmd)
}

// GetConnectionCount returns the number of active connections to other peers.
func (c *Client) GetConnectionCount() (count int64, err error) {
	return c.GetConnectionCountAsync().Receive()
}

// FuturePingResult is a future promise to deliver the result of a PingAsync RPC invocation (or an applicable error).
type FuturePingResult chan *response

// Receive waits for the response promised by the future and returns the result of queueing a ping to be sent to each connected peer.
func (r FuturePingResult) Receive() (err error) {
	_, err = receiveFuture(r)
	return err
}

// PingAsync returns an instance of a type that can be used to get the result of the RPC at some future time by invoking the Receive function on the returned instance. See Ping for the blocking version and more details.
func (c *Client) PingAsync() FuturePingResult {
	cmd := btcjson.NewPingCmd()
	return c.sendCmd(cmd)
}

// Ping queues a ping to be sent to each connected peer. Use the GetPeerInfo function and examine the PingTime and PingWait fields to access the ping times.
func (c *Client) Ping() (err error) {
	return c.PingAsync().Receive()
}

// FutureGetPeerInfoResult is a future promise to deliver the result of a GetPeerInfoAsync RPC invocation (or an applicable error).
type FutureGetPeerInfoResult chan *response

// Receive waits for the response promised by the future and returns  data about each connected network peer.
func (r FutureGetPeerInfoResult) Receive() (peerInfo []btcjson.GetPeerInfoResult, err error) {
	var res []byte
	if res, err = receiveFuture(r); slog.Check(err) {
		return
	}
	// Unmarshal result as an array of getpeerinfo result objects.
	if err = js.Unmarshal(res, &peerInfo); slog.Check(err) {
		return
	}
	return
}

// GetPeerInfoAsync returns an instance of a type that can be used to get the result of the RPC at some future time by invoking the Receive function on the returned instance. See GetPeerInfo for the blocking version and more details.
func (c *Client) GetPeerInfoAsync() FutureGetPeerInfoResult {
	cmd := btcjson.NewGetPeerInfoCmd()
	return c.sendCmd(cmd)
}

// GetPeerInfo returns data about each connected network peer.
func (c *Client) GetPeerInfo() (res []btcjson.GetPeerInfoResult, err error) {
	return c.GetPeerInfoAsync().Receive()
}

// FutureGetNetTotalsResult is a future promise to deliver the result of a GetNetTotalsAsync RPC invocation (or an applicable error).
type FutureGetNetTotalsResult chan *response

// Receive waits for the response promised by the future and returns network statistics.
func (r FutureGetNetTotalsResult) Receive() (totals *btcjson.GetNetTotalsResult, err error) {
	res, err := receiveFuture(r)
	if err != nil {
		slog.Error(err)
		return nil, err
	}
	// Unmarshal result as a getnettotals result object.
	totals = &btcjson.GetNetTotalsResult{}
	if err = js.Unmarshal(res, &totals); slog.Check(err) {
		return
	}
	return
}

// GetNetTotalsAsync returns an instance of a type that can be used to get the result of the RPC at some future time by invoking the Receive function on the returned instance. See GetNetTotals for the blocking version and more details.
func (c *Client) GetNetTotalsAsync() FutureGetNetTotalsResult {
	cmd := btcjson.NewGetNetTotalsCmd()
	return c.sendCmd(cmd)
}

// GetNetTotals returns network traffic statistics.
func (c *Client) GetNetTotals() (res *btcjson.GetNetTotalsResult, err error) {
	return c.GetNetTotalsAsync().Receive()
}
