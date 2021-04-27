package checkpoints

import (
	"fmt"
	"github.com/p9c/matrjoska/pkg/chaincfg"
	"github.com/p9c/matrjoska/pkg/chainhash"
	"strconv"
	"strings"
)

// Parse checks the checkpoint strings for valid syntax ( '<height>:<hash>') and parses them to
// chaincfg.Checkpoint instances.
func Parse(checkpointStrings []string) ([]chaincfg.Checkpoint, error) {
	if len(checkpointStrings) == 0 {
		return nil, nil
	}
	checkpoints := make([]chaincfg.Checkpoint, len(checkpointStrings))
	for i, cpString := range checkpointStrings {
		checkpoint, e := newCheckpointFromStr(cpString)
		if e != nil {
			return nil, e
		}
		checkpoints[i] = checkpoint
	}
	return checkpoints, nil
}

// newCheckpointFromStr parses checkpoints in the '<height>:<hash>' format.
func newCheckpointFromStr(checkpoint string) (chaincfg.Checkpoint, error) {
	parts := strings.Split(checkpoint, ":")
	if len(parts) != 2 {
		return chaincfg.Checkpoint{}, fmt.Errorf(
			"unable to parse "+
				"checkpoint %q -- use the syntax <height>:<hash>",
			checkpoint,
		)
	}
	height, e := strconv.ParseInt(parts[0], 10, 32)
	if e != nil {
		return chaincfg.Checkpoint{}, fmt.Errorf(
			"unable to parse "+
				"checkpoint %q due to malformed height", checkpoint,
		)
	}
	if len(parts[1]) == 0 {
		return chaincfg.Checkpoint{}, fmt.Errorf(
			"unable to parse "+
				"checkpoint %q due to missing hash", checkpoint,
		)
	}
	hash, e := chainhash.NewHashFromStr(parts[1])
	if e != nil {
		return chaincfg.Checkpoint{}, fmt.Errorf(
			"unable to parse "+
				"checkpoint %q due to malformed hash", checkpoint,
		)
	}
	return chaincfg.Checkpoint{
			Height: int32(height),
			Hash:   hash,
		},
		nil
}
