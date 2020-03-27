package main

import (
	"fmt"
	"github.com/p9c/pod/pkg/pipe"
	"time"
)

func main() {
	p := pipe.Parent(func(b []byte) (err error) {
		fmt.Print("from parent: ", string(b))
		return
	}, make(chan struct{}))
	for {
		p.Write([]byte("ping"))
		time.Sleep(time.Second)
	}
}
