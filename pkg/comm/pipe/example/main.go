package main

import (
	"fmt"
	"time"
	
	qu "github.com/p9c/pod/pkg/util/quit"
	
	"github.com/p9c/pod/pkg/comm/pipe"
)

func main() {
	quit := qu.T()
	p := pipe.Consume(
		quit, func(b []byte) (e error) {
			fmt.Println("from child:", string(b))
			return
		}, "go", "run", "serve/main.go",
	)
	for {
		_, e := p.StdConn.Write([]byte("ping"))
		if e != nil  {
			fmt.Println("err:", err)
		}
		time.Sleep(time.Second)
	}
}
