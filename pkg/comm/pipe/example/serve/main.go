package main

import (
	"fmt"
	qu "github.com/p9c/pod/pkg/util/quit"
	"time"

	"github.com/p9c/pod/pkg/comm/pipe"
)

func main() {
	p := pipe.Serve(qu.T(), func(b []byte) (err error) {
		fmt.Print("from parent: ", string(b))
		return
	})
	for {
		_, err := p.Write([]byte("ping"))
		if err != nil {
			fmt.Println("err:", err)
		}
		time.Sleep(time.Second)
	}
}
