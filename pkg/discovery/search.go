package discovery

import (
	"context"
	"fmt"
	"os"

	"github.com/grandcat/zeroconf"

	"git.parallelcoin.io/dev/pod/pkg/util/cl"
)

type ResultsChan chan *zeroconf.ServiceEntry

func AsyncZeroConfSearch(service, group string) (cancel context.CancelFunc,
	r ResultsChan, err error) {
	r = make(ResultsChan, 10)
	myInstance := fmt.Sprint(os.Getppid())
	domain := "local."
	log <- cl.Warn{"starting search", cl.Ine()}
	resolver, err := zeroconf.NewResolver(nil)
	if err != nil {
		log <- cl.Error{"Failed to initialize resolver:", err.Error()}
		return
	}
	entries := make(chan *zeroconf.ServiceEntry)
	go func(results <-chan *zeroconf.ServiceEntry) {
		for entry := range results {
			if entry.Service == service &&
				len(entry.Text) > 0 &&
				entry.Text[0] == "group="+group {
				if entry.Instance == myInstance {
					continue
				} else {
					r <-entry
				}
			}
		}
	}(entries)
	var ctx context.Context
	ctx, cancel = context.WithCancel(context.Background())
	err = resolver.Browse(ctx, service, domain, entries)
	if err != nil {
		log <- cl.Error{"Failed to browse:", err.Error()}
		panic(err)
	}
	return
}
