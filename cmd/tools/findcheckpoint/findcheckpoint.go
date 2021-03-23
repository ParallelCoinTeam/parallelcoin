package main

import (
	"fmt"
	"path/filepath"
	
	"github.com/p9c/pod/pkg/blockchain"
	"github.com/p9c/pod/pkg/chaincfg"
	"github.com/p9c/pod/pkg/chainhash"
	"github.com/p9c/pod/pkg/database"
)

const blockDbNamePrefix = "blocks"

var (
	cfg *config
)

// loadBlockDB opens the block database and returns a handle to it.
func loadBlockDB() (database.DB, error) {
	// The database name is based on the database type.
	dbName := blockDbNamePrefix + "_" + cfg.DbType
	dbPath := filepath.Join(cfg.DataDir, dbName)
	I.F("Loading block database from '%s'\n", dbPath)
	db, e := database.Open(cfg.DbType, dbPath, activeNetParams.Net)
	if e != nil {
		return nil, e
	}
	return db, nil
}

// findCandidates searches the chain backwards for checkpoint candidates and
// returns a slice of found candidates,
// if any.  It also stops searching for candidates at the last checkpoint
// that is already hard coded into btcchain since there is no point in
// finding candidates before already existing checkpoints.
func findCandidates(
	chain *blockchain.BlockChain, latestHash *chainhash.Hash,
) ([]*chaincfg.Checkpoint, error) {
	// Start with the latest block of the main chain.
	block, e := chain.BlockByHash(latestHash)
	if e != nil {
		return nil, e
	}
	// Get the latest known checkpoint.
	latestCheckpoint := chain.LatestCheckpoint()
	if latestCheckpoint == nil {
		// Set the latest checkpoint to the genesis block if there isn't
		// already one.
		latestCheckpoint = &chaincfg.Checkpoint{
			Hash:   activeNetParams.GenesisHash,
			Height: 0,
		}
	}
	// The latest known block must be at least the last known checkpoint plus
	// required checkpoint confirmations.
	checkpointConfirmations := int32(blockchain.CheckpointConfirmations)
	requiredHeight := latestCheckpoint.Height + checkpointConfirmations
	if block.Height() < requiredHeight {
		return nil, fmt.Errorf(
			"the block database is only at height "+
				"%d which is less than the latest checkpoint height "+
				"of %d plus required confirmations of %d",
			block.Height(), latestCheckpoint.Height,
			checkpointConfirmations,
		)
	}
	// For the first checkpoint,
	// the required height is any block after the genesis block,
	// so long as the chain has at least the required number of confirmations
	// (which is enforced above).
	if len(activeNetParams.Checkpoints) == 0 {
		requiredHeight = 1
	}
	// Indeterminate progress setup.
	numBlocksToTest := block.Height() - requiredHeight
	progressInterval := (numBlocksToTest / 100) + 1 // min 1
	fmt.Print("Searching for candidates")
	defer fmt.Println()
	// Loop backwards through the chain to find checkpoint candidates.
	candidates := make([]*chaincfg.Checkpoint, 0, cfg.NumCandidates)
	numTested := int32(0)
	for len(candidates) < cfg.NumCandidates && block.Height() > requiredHeight {
		// Display progress.
		if numTested%progressInterval == 0 {
			fmt.Print(".")
		}
		// Determine if this block is a checkpoint candidate.
		isCandidate, e := chain.IsCheckpointCandidate(block)
		if e != nil {
			return nil, e
		}
		// All checks passed, so this node seems like a reasonable checkpoint candidate.
		if isCandidate {
			checkpoint := chaincfg.Checkpoint{
				Height: block.Height(),
				Hash:   block.Hash(),
			}
			candidates = append(candidates, &checkpoint)
		}
		prevHash := &block.MsgBlock().Header.PrevBlock
		block, e = chain.BlockByHash(prevHash)
		if e != nil {
			return nil, e
		}
		numTested++
	}
	return candidates, nil
}

// showCandidate display a checkpoint candidate using and output format determined by the configuration parameters.  The Go syntax output uses the format the btcchain code expects for checkpoints added to the list.
func showCandidate(
	candidateNum int, checkpoint *chaincfg.Checkpoint,
) {
	if cfg.UseGoOutput {
		I.F(
			"Candidate %d -- {%d, newShaHashFromStr(\"%v\")},\n",
			candidateNum, checkpoint.Height, checkpoint.Hash,
		)
		return
	}
	I.F(
		"Candidate %d -- Height: %d, Hash: %v\n", candidateNum,
		checkpoint.Height, checkpoint.Hash,
	)
}
func main() {
	// Load configuration and parse command line.
	var tcfg *config
	var e error
	tcfg, _, e = loadConfig()
	if e != nil {
		return
	}
	cfg = tcfg
	// Load the block database.
	db, e := loadBlockDB()
	if e != nil {
		E.Ln("failed to load database:", e)
		return
	}
	defer func() {
		if e := db.Close(); E.Chk(e) {
		}
	}()
	// Setup chain.  Ignore notifications since they aren't needed for this util.
	chain, e := blockchain.New(
		&blockchain.Config{
			DB:          db,
			ChainParams: activeNetParams,
			TimeSource:  blockchain.NewMedianTime(),
		},
	)
	if e != nil {
		E.Ln("failed to initialize chain: %v\n", e)
		return
	}
	// Get the latest block hash and height from the database and report status.
	best := chain.BestSnapshot()
	I.F("Block database loaded with block height %d\n", best.Height)
	// Find checkpoint candidates.
	candidates, e := findCandidates(chain, &best.Hash)
	if e != nil {
		E.Ln("Unable to identify candidates:", e)
		return
	}
	// No candidates.
	if len(candidates) == 0 {
		E.Ln("No candidates found.")
		return
	}
	// Show the candidates.
	for i, checkpoint := range candidates {
		showCandidate(i+1, checkpoint)
	}
}
