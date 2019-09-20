package main

import (
	"fmt"

	"github.com/parallelcointeam/parallelcoin/pkg/util/interrupt"
)

func main() {
	interrupt.AddHandler(func() {
		fmt.Println("IT'S THE END OF THE WORLD!")
	})
	<-interrupt.HandlersDone
}
