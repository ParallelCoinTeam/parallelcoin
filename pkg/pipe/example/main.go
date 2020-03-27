package main

import (
	"fmt"
	"github.com/p9c/pod/pkg/pipe"
	"time"
)

func main() {
	p := pipe.Child(func(b []byte) (err error) {
		fmt.Println("from child:", string(b))
		return
	}, "go", "run", "serve/main.go")
	for {
		_, err := p.Write([]byte("ping"))
		if err != nil {
			fmt.Println("err:", err)
		}
		time.Sleep(time.Second)
	}
}
