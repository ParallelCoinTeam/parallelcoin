package main

import (
	"github.com/p9c/pod/pkg/log"
	"github.com/p9c/pod/pkg/util/interrupt"
)

func main() {
	interrupt.AddHandler(func() {
		log.Println("IT'S THE END OF THE WORLD!")
	})
	<-interrupt.HandlersDone
}
