package main

import (
	"github.com/p9c/pkg/app/slog"
	"log"

	client "github.com/p9c/pod/pkg/rpc/client"
)

func main() {
	// Connect to local bitcoin core RPC server using HTTP POST mode.
	connCfg := &client.ConnConfig{
		Host:         "localhost:11046",
		User:         "yourrpcuser",
		Pass:         "yourrpcpass",
		HTTPPostMode: true,  // Bitcoin core only supports HTTP POST mode
		TLS:          false, // Bitcoin core does not provide TLS by default
	}
	// Notice the notification parameter is nil since notifications are not supported in HTTP POST mode.
	cl, err := client.New(connCfg, nil)
	if err != nil {
		slog.Fatal(err)
	}
	defer cl.Shutdown()
	// Get the current block count.
	blockCount, err := cl.GetBlockCount()
	if err != nil {
		slog.Fatal(err)
	}
	log.Printf("Block count: %d", blockCount)
}
