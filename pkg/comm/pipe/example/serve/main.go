package main

import (
	"fmt"
	"time"
	
	qu "github.com/p9c/pod/pkg/util/qu"
	
	"github.com/p9c/pod/pkg/comm/pipe"
)

func main() {
	p := pipe.Serve(qu.T(), func(b []byte) (e error) {
		fmt.Print("from parent: ", string(b))
		return
	})
	for {
		_, e := p.Write([]byte("ping"))
		if e != nil  {
			fmt.Println("err:", err)
		}
		time.Sleep(time.Second)
	}
}
