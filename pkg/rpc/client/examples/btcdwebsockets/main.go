package main

import (
	"io/ioutil"
	"log"
	"path/filepath"
	"time"

	"github.com/p9c/pod/app/appdata"
	"github.com/p9c/pod/pkg/chain/wire"
	rpcclient "github.com/p9c/pod/pkg/rpc/client"
	"github.com/p9c/pod/pkg/util"
)

func main() {
	// Only override the handlers for notifications you care about. Also note most of these handlers will only be called if you register for notifications.  See the documentation of the rpcclient NotificationHandlers type for more details about each handler.
	ntfnHandlers := rpcclient.NotificationHandlers{
		OnFilteredBlockConnected: func(height int32, header *wire.BlockHeader, txns []*util.Tx) {
			log.Printf("Block connected: %v (%d) %v",
				header.BlockHash(), height, header.Timestamp)
		},
		OnFilteredBlockDisconnected: func(height int32, header *wire.BlockHeader) {
			log.Printf("Block disconnected: %v (%d) %v",
				header.BlockHash(), height, header.Timestamp)
		},
	}
	// Connect to local pod RPC server using websockets.
	podHomeDir := appdata.Dir("pod", false)
	certs, err := ioutil.ReadFile(filepath.Join(podHomeDir, "rpc.cert"))
	if err != nil {
		Fatal(err)
	}
	connCfg := &rpcclient.ConnConfig{
		Host:         "localhost:11048",
		Endpoint:     "ws",
		User:         "yourrpcuser",
		Pass:         "yourrpcpass",
		Certificates: certs,
	}
	client, err := rpcclient.New(connCfg, &ntfnHandlers)
	if err != nil {
		Fatal(err)
	}
	// Register for block connect and disconnect notifications.
	if err := client.NotifyBlocks(); err != nil {
		Fatal(err)
	}
	fmt.Println("NotifyBlocks: Registration Complete")
	// Get the current block count.
	blockCount, err := client.GetBlockCount()
	if err != nil {
		Fatal(err)
	}
	log.Printf("Block count: %d", blockCount)
	// For this example gracefully shutdown the client after 10 seconds. Ordinarily when to shutdown the client is highly application specific.
	fmt.Println("Client shutdown in 10 seconds...")
	time.AfterFunc(time.Second*10, func() {
		fmt.Println("Client shutting down...")
		client.Shutdown()
		fmt.Println("Client shutdown complete.")
	})
	// Wait until the client either shuts down gracefully (or the user terminates the process with Ctrl+C).
	client.WaitForShutdown()
}
