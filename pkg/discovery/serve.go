package discovery

import (
	"context"
	"fmt"
	"net"
	"os"
	"strings"

	"github.com/grandcat/zeroconf"

	"github.com/parallelcointeam/parallelcoin/pkg/chain/config/netparams"
)

type Request struct {
	Key     string
	Address string
}

type RequestFunc func(key, address string)

func GetParallelcoinServiceName(params *netparams.Params) string {
	return fmt.Sprintf("parallelcoin/%s", params.Net)
}

func Serve(params *netparams.Params, lanInterface *net.Interface,
	group string) (cancel context.CancelFunc, request RequestFunc, err error) {
	WARN("starting discovery server")
	texts := []string{"group=" + group}
	domain := "local."
	requests := make(chan Request)
	request = func(key, address string) {
		requests <- Request{key, address}
	}
	ctx, cancel := context.WithCancel(context.Background())
	go func() {
		alias := fmt.Sprint(os.Getppid())
		server, err := zeroconf.Register(alias, GetParallelcoinServiceName(
			params), domain, 1, texts, []net.Interface{*lanInterface})
		if err != nil {
			ERROR("error registering ", err)
			return
		}
		TRACE("registered")
		for {
			select {
			case r := <-requests:
				found := false
				for i := range texts {
					split := strings.Split(texts[i], "=")
					// WARN("'",split[0],"' '", r.Key,"'")
					if split[0] == r.Key {
						found = true
						if r.Address == "" {
							DEBUG("discovery: removing key ", r.Key)
							switch {
							case i == 0:
								texts = texts[1:]
							case i == len(texts)-1:
								texts = texts[:len(texts)-1]
							default:
								texts = append(texts[:i], texts[i+1:]...)
							}
						} else {
							texts[i] = r.Key + "=" + r.Address
						}
						break
					}
				}
				if !found && r.Address != "" {
					nt := r.Key + "=" + r.Address
					texts = append(texts, nt)
					DEBUG("appending ", nt, " to texts ", texts)
				}
				server.Shutdown()
				TRACE("shut down server")
				server, err = zeroconf.Register(alias,
					GetParallelcoinServiceName(
						params), domain, 1, texts, []net.Interface{*lanInterface})
				if err != nil {
					ERROR("error registering ", err)
					return
				}
				TRACE("restarted server")
			case <-ctx.Done():
				server.Shutdown()
				break
			}
		}
	}()
	return
}
