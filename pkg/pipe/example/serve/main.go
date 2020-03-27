package main

import (
	"fmt"
	"github.com/p9c/pod/pkg/pipe"
	"time"
)

func main() {
	p := pipe.Serve(func(b []byte) (err error) {
		fmt.Print("from parent: ", string(b))
		return
	}, make(chan struct{}))
	for {
		_, err := p.Write([]byte("ping"))
		if err != nil {
			fmt.Println("err:", err)
		}
		time.Sleep(time.Second)
	}
}
