package main

import (
	"fmt"
	"github.com/stalker-loki/app/slog"
	"io/ioutil"
	"log"
	"path/filepath"
	"time"

	"github.com/davecgh/go-spew/spew"

	"github.com/stalker-loki/pod/app/appdata"
	client "github.com/stalker-loki/pod/pkg/rpc/client"
	"github.com/stalker-loki/pod/pkg/util"
)

func main() {
	// Only override the handlers for notifications you care about. Also note most of the handlers will only be called if you register for notifications.  See the documentation of the cl NotificationHandlers type for more details about each handler.
	ntfnHandlers := client.NotificationHandlers{
		OnAccountBalance: func(account string, balance util.Amount, confirmed bool) {
			log.Printf("New balance for account %s: %v", account,
				balance)
		},
	}
	// Connect to local btcwallet RPC server using websockets.
	certHomeDir := appdata.Dir("mod", false)
	certs, err := ioutil.ReadFile(filepath.Join(certHomeDir, "rpc.cert"))
	if err != nil {
		slog.Fatal(err)
	}
	connCfg := &client.ConnConfig{
		Host:         "localhost:11046",
		Endpoint:     "ws",
		User:         "yourrpcuser",
		Pass:         "yourrpcpass",
		Certificates: certs,
	}
	cl, err := client.New(connCfg, &ntfnHandlers)
	if err != nil {
		slog.Fatal(err)
	}
	// Get the list of unspent transaction outputs (utxos) that the connected wallet has at least one private key for.
	unspent, err := cl.ListUnspent()
	if err != nil {
		slog.Fatal(err)
	}
	log.Printf("Num unspent outputs (utxos): %d", len(unspent))
	if len(unspent) > 0 {
		log.Printf("First utxo:\n%v", spew.Sdump(unspent[0]))
	}
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
