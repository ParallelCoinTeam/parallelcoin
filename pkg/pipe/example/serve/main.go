package main

import (
	"fmt"
	"io"
	"os"
	"time"
)

func main() {
	var n int
	var err error
	data := make([]byte, 1024)
	go func() {
		for {
			n, err = os.Stdin.Read(data)
			if n > 0 {
				fmt.Print("from parent: ", string(data[:n]))
			}

			if err != nil && err != io.EOF {
				fmt.Println("err: ", err)
			}
		}
	}()
	for {
		io.WriteString(os.Stdout, "ping")
		time.Sleep(time.Second)
	}
}
