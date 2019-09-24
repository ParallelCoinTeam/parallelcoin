// +build ignore

package main

import (
	"fmt"
	"os"
	"time"

	"github.com/davecgh/go-spew/spew"

	"github.com/parallelcointeam/parallelcoin/pkg/discovery"
	"github.com/parallelcointeam/parallelcoin/pkg/chain/config/netparams"
)

func main() {
	routeable := discovery.GetRouteableInterface()
	stopServe, request, err := discovery.Serve(&netparams.TestNet3Params, routeable)
	if err != nil {
		fmt.Println(err)
		panic(err)
	}
	fmt.Println("started  autoconf server")
	service := discovery.GetParallelcoinServiceName(&netparams.
		TestNet3Params)
	stopSearch, results, err := discovery.AsyncZeroConfSearch(service, fmt.Sprint(os.Getppid()))
	if err != nil {
		fmt.Println(err)
		panic(err)
	}
	time.Sleep(time.Second)
	r := results()
	fmt.Println("default initial")
	for i := range r {
		fmt.Println("current", r[i].Instance, r[i].Text)
	}
	fmt.Println("adding txt entry")
	request("testing", "1.1.1.1")
	time.Sleep(time.Second)
	r = results()
	for i := range r {
		fmt.Println("current",i, r[i].Instance, r[i].Text)
		spew.Dump(r[i])
	}
	fmt.Println("removing txt entry")
	request("testing", "")
	time.Sleep(time.Second)
	r = results()
	for i := range r {
		fmt.Println("current", r[i].Instance, r[i].Text)
	}
	stopSearch()
	stopServe()
}
