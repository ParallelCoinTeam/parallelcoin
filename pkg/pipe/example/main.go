package main

import (
	"fmt"
	"github.com/p9c/pod/pkg/stdconn/worker"
	"io"
	"time"
)

func main() {
	var n int
	var err error
	w := worker.Spawn("go", "run", "serve/main.go")
	data := make([]byte, 1024)
	go func() {
		for {
			n, err = w.StdConn.Read(data)
			if n > 0 {
				fmt.Println("from child:", string(data[:n]))
			}

			if err != nil && err != io.EOF {
				fmt.Println("err:",err)
			}
		}
	}()
	for {
		_, err = w.StdConn.Write([]byte("ping"))
		if err != nil {
			fmt.Println("err:", err)
		}
		time.Sleep(time.Second)
	}
}
