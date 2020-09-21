package main

import (
	"fmt"
	"github.com/stalker-loki/pod/pkg/util/interrupt"
)

func main() {
	interrupt.AddHandler(func() {
		fmt.Println("IT'S THE END OF THE WORLD!")
	})
	<-interrupt.HandlersDone
}
