package main

import (
	log "github.com/p9c/pod/pkg/logi"
	"github.com/p9c/pod/pkg/util/interrupt"
)

func main() {
	interrupt.AddHandler(func() {
		log.Println("IT'S THE END OF THE WORLD!")
	})
	<-interrupt.HandlersDone
}
