package rcd

import (
	blockchain "github.com/p9c/pod/pkg/chain"
)

const (
	NewBlock uint32 = iota
	// Add new events here
	EventCount
)

type Event struct {
	Type    uint32
	Payload []byte
}

var EventsChan = make(chan Event, 24)

func ListenInit(rc RcVar) {
	rc.Events = EventsChan
	rc.Cx.RealNode.Chain.Subscribe(func(callback *blockchain.Notification) {
		switch callback.Type {
		case blockchain.NTBlockAccepted:
			//callback.
		}
	})

	return
}
