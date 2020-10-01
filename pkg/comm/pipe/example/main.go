package main

import (
	"fmt"
	"time"

	"github.com/p9c/pod/pkg/comm/pipe"
)

func main() {
	quit := make(chan struct{})
	p := pipe.Consume(quit, func(b []byte) (err error) {
		fmt.Println("from child:", string(b))
		return
	}, "go", "run", "serve/main.go")
	for {
		_, err := p.StdConn.Write([]byte("ping"))
		if err != nil {
			fmt.Println("err:", err)
		}
		time.Sleep(time.Second)
	}
}
