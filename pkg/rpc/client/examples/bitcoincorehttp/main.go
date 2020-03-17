package main

import (
	"log"

	rpcclient "github.com/p9c/pod/pkg/rpc/client"
)

func main() {
	// Connect to local bitcoin core RPC server using HTTP POST mode.
	connCfg := &rpcclient.ConnConfig{
		Host:         "localhost:11046",
		User:         "yourrpcuser",
		Pass:         "yourrpcpass",
		HTTPPostMode: true,  // Bitcoin core only supports HTTP POST mode
		TLS:          false, // Bitcoin core does not provide TLS by default
	}
	// Notice the notification parameter is nil since notifications are not supported in HTTP POST mode.
	client, err := rpcclient.New(connCfg, nil)
	if err != nil {
		L.Fatal(err)
	}
	defer client.Shutdown()
	// Get the current block count.
	blockCount, err := client.GetBlockCount()
	if err != nil {
		L.Fatal(err)
	}
	log.Printf("Block count: %d", blockCount)
}
