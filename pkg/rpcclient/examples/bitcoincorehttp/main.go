package main

import (
	"log"
	
	qu "github.com/p9c/qu"
	
	rpcclient "github.com/p9c/pod/pkg/rpcclient"
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
	client, e := rpcclient.New(connCfg, nil, qu.T())
	if e != nil  {
		F.Ln(e)
	}
	defer client.Shutdown()
	// Get the current block count.
	blockCount, e := client.GetBlockCount()
	if e != nil  {
		F.Ln(e)
	}
	log.Printf("Block count: %d", blockCount)
}
