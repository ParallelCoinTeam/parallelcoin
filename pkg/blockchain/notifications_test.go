package blockchain

import (
	"testing"
	
	"github.com/p9c/matrjoska/pkg/chaincfg"
)

// TestNotifications ensures that notification callbacks are fired on events.
func TestNotifications(t *testing.T) {
	blocks, e := loadBlocks("blk_0_to_4.dat.bz2")
	if e != nil {
		t.Fatalf("Error loading file: %v\n", e)
	}
	// Create a new database and chain instance to run tests against.
	chain, teardownFunc, e := chainSetup("notifications",
		&chaincfg.MainNetParams,
	)
	if e != nil {
		t.Fatalf("Failed to setup chain instance: %v", e)
	}
	defer teardownFunc()
	notificationCount := 0
	callback := func(notification *Notification) {
		if notification.Type == NTBlockAccepted {
			notificationCount++
		}
	}
	// Register callback multiple times then assert it is called that many times.
	const numSubscribers = 3
	for i := 0; i < numSubscribers; i++ {
		chain.Subscribe(callback)
	}
	_, _, e = chain.ProcessBlock(0, blocks[1], BFNone, blocks[1].Height())
	if e != nil {
		t.Fatalf("ProcessBlock fail on block 1: %v\n", e)
	}
	if notificationCount != numSubscribers {
		t.Fatalf("Expected notification callback to be executed %d "+
			"times, found %d", numSubscribers, notificationCount,
		)
	}
}
