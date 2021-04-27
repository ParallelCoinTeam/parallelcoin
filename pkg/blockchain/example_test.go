package blockchain_test

import (
	"fmt"
	bits2 "github.com/p9c/matrjoska/pkg/bits"
	"github.com/p9c/matrjoska/pkg/block"
	"log"
	"math/big"
	"os"
	"path/filepath"
	
	"github.com/p9c/matrjoska/pkg/blockchain"
	"github.com/p9c/matrjoska/pkg/chaincfg"
	"github.com/p9c/matrjoska/pkg/database"
	_ "github.com/p9c/matrjoska/pkg/database/ffldb"
)

// This example demonstrates how to create a new chain instance and use ProcessBlock to attempt to add a block to the
// chain. As the package overview documentation describes, this includes all of the Bitcoin consensus rules. This
// example intentionally attempts to insert a duplicate genesis block to illustrate how an invalid block is handled.
func ExampleBlockChain_ProcessBlock() {
	// Create a new database to store the accepted blocks into. Typically this would be opening an existing database and
	// would not be deleting and creating a new database like this, but it is done here so this is a complete working
	// example and does not leave temporary files laying around.
	dbPath := filepath.Join(os.TempDir(), "exampleprocessblock")
	_ = os.RemoveAll(dbPath)
	db, e := database.Create("ffldb", dbPath, chaincfg.MainNetParams.Net)
	if e != nil {
		log.Printf("Failed to create database: %v\n", e)
		return
	}
	defer func() {
		if e = os.RemoveAll(dbPath); E.Chk(e) {
		}
	}()
	defer func() {
		if e = db.Close(); E.Chk(e) {
		}
	}()
	// Create a new BlockChain instance using the underlying database for the main bitcoin network. This example does
	// not demonstrate some of the other available configuration options such as specifying a notification callback and
	// signature cache. Also, the caller would ordinarily keep a reference to the median time source and add time values
	// obtained from other peers on the network so the local time is adjusted to be in agreement with other peers.
	chain, e := blockchain.New(
		&blockchain.Config{
			DB:          db,
			ChainParams: &chaincfg.MainNetParams,
			TimeSource:  blockchain.NewMedianTime(),
		},
	)
	if e != nil {
		log.Printf("Failed to create chain instance: %v\n", e)
		return
	}
	// Process a block. For this example, we are going to intentionally cause an error by trying to process the genesis
	// block which already exists.
	genesisBlock := block.NewBlock(chaincfg.MainNetParams.GenesisBlock)
	var isMainChain bool
	var isOrphan bool
	isMainChain, isOrphan, e = chain.ProcessBlock(
		0, genesisBlock,
		blockchain.BFNone, 0,
	)
	if e != nil {
		log.Printf("Failed to process block: %v\n", e)
		return
	}
	log.Printf("Block accepted. Is it on the main chain?: %v", isMainChain)
	log.Printf("Block accepted. Is it an orphan?: %v", isOrphan)
	// Output:
	// Failed to process block: already have block 000000000019d6689c085ae165831e934ff763ae46a2a6c172b3f1b60a8ce26f
}

// This example demonstrates how to convert the compact "bits" in a block header which represent the target difficulty
// to a big integer and display it using the typical hex notation.
func ExampleCompactToBig() {
	// Convert the bits from block 300000 in the main block chain.
	bits := uint32(419465580)
	targetDifficulty := bits2.CompactToBig(bits)
	// Display it in hex.
	fmt.Printf("%064x\n", targetDifficulty.Bytes())
	// Output:
	// 0000000000000000896c00000000000000000000000000000000000000000000
}

// This example demonstrates how to convert a target difficulty into the compact "bits" in a block header which
// represent that target difficulty .
func ExampleBigToCompact() {
	// Convert the target difficulty from block 300000 in the main block
	// chain to compact form.
	t := "0000000000000000896c00000000000000000000000000000000000000000000"
	targetDifficulty, success := new(big.Int).SetString(t, 16)
	if !success {
		fmt.Println("invalid target difficulty")
		return
	}
	bits := bits2.BigToCompact(targetDifficulty)
	fmt.Println(bits)
	// Output:
	// 419465580
}
