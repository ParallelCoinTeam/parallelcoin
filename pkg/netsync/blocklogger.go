package netsync

import (
	"fmt"
	"github.com/p9c/matrjoska/pkg/block"
	"sync"
	"time"
)

// blockProgressLogger provides periodic logging for other services in order to show users progress of certain "actions"
// involving some or all current blocks. Ex: syncing to best chain, indexing all blocks, etc.
type blockProgressLogger struct {
	receivedLogBlocks int64
	receivedLogTx     int64
	lastBlockLogTime  time.Time
	// subsystemLogger   *log.Logger
	progressAction string
	sync.Mutex
}

// newBlockProgressLogger returns a new block progress logger. The progress message is templated as follows:
// {progressAction } {numProcessed} {blocks|block} in the last {timePeriod} ({numTxs}, height {lastBlockHeight},
// {lastBlockTimeStamp})
func newBlockProgressLogger(progressMessage string) *blockProgressLogger {
	return &blockProgressLogger{
		lastBlockLogTime: time.Now(),
		progressAction:   progressMessage,
	}
}

// LogBlockHeight logs a new block height as an information message to show progress to the user.
//
// In order to prevent spam, it limits logging to one message every 10 seconds with duration and totals included.
func (b *blockProgressLogger) LogBlockHeight(block *block.Block) {
	b.Lock()
	defer b.Unlock()
	b.receivedLogBlocks++
	b.receivedLogTx += int64(len(block.WireBlock().Transactions))
	now := time.Now()
	duration := now.Sub(b.lastBlockLogTime)
	if duration < time.Second*2 {
		return
	}
	// Truncate the duration to 10s of milliseconds.
	durationMillis := int64(duration / time.Millisecond)
	tDuration := 10 * time.Millisecond * time.Duration(durationMillis/10)
	// Log information about new block height.
	blockStr := "blocks"
	if b.receivedLogBlocks == 1 {
		blockStr = "block "
	}
	txStr := "transactions"
	if b.receivedLogTx == 1 {
		txStr = "transaction "
	}
	tD := tDuration.Seconds()
	I.F(
		"%s %6d %s in the last %s (%6d %s, height %8d, %s) %0.2f tx/s",
		b.progressAction,
		b.receivedLogBlocks,
		blockStr,
		fmt.Sprintf("%0.1fs", tD),
		b.receivedLogTx,
		txStr, block.Height(),
		block.WireBlock().Header.Timestamp,
		float64(b.receivedLogTx)/tD,
	)
	b.receivedLogBlocks = 0
	b.receivedLogTx = 0
	b.lastBlockLogTime = now
}

func (b *blockProgressLogger) SetLastLogTime(time time.Time) {
	b.lastBlockLogTime = time
}
