package main

import (
	"fmt"
	"github.com/stalker-loki/app/slog"
	"io/ioutil"
	"log"
	"path/filepath"
	"time"

	"github.com/stalker-loki/pod/app/appdata"
	"github.com/stalker-loki/pod/pkg/chain/wire"
	client "github.com/stalker-loki/pod/pkg/rpc/client"
	"github.com/stalker-loki/pod/pkg/util"
)

func main() {
	// Only override the handlers for notifications you care about. Also note most of these handlers will only be called
	// if you register for notifications.  See the documentation of the cl NotificationHandlers type for more
	// details about each handler.
	ntfnHandlers := client.NotificationHandlers{
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
		slog.Fatal(err)
	}
	connCfg := &client.ConnConfig{
		Host:         "localhost:11048",
		Endpoint:     "ws",
		User:         "yourrpcuser",
		Pass:         "yourrpcpass",
		Certificates: certs,
	}
	cl, err := client.New(connCfg, &ntfnHandlers)
	if err != nil {
		slog.Fatal(err)
	}
	// Register for block connect and disconnect notifications.
	if err := cl.NotifyBlocks(); err != nil {
		slog.Fatal(err)
	}
	fmt.Println("NotifyBlocks: Registration Complete")
	// Get the current block count.
	blockCount, err := cl.GetBlockCount()
	if err != nil {
		slog.Fatal(err)
	}
	log.Printf("Block count: %d", blockCount)
	// For this example gracefully shutdown the cl after 10 seconds. Ordinarily when to shutdown the cl is highly application specific.
	fmt.Println("Client shutdown in 10 seconds...")
	time.AfterFunc(time.Second*10, func() {
		fmt.Println("Client shutting down...")
		cl.Shutdown()
		fmt.Println("Client shutdown complete.")
	})
	// Wait until the cl either shuts down gracefully (or the user terminates the process with Ctrl+C).
	cl.WaitForShutdown()
}
