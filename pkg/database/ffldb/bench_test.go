package ffldb

import (
	block2 "github.com/p9c/matrjoska/pkg/block"
	"os"
	"path/filepath"
	"testing"
	
	"github.com/p9c/matrjoska/pkg/chaincfg"
	"github.com/p9c/matrjoska/pkg/database"
)

// BenchmarkBlockHeader benchmarks how long it takes to load the mainnet genesis block header.
func BenchmarkBlockHeader(b *testing.B) {
	// Start by creating a new database and populating it with the mainnet genesis block.
	dbPath := filepath.Join(os.TempDir(), "ffldb-benchblkhdr")
	_ = os.RemoveAll(dbPath)
	db, e := database.Create("ffldb", dbPath, blockDataNet)
	if e != nil {
		b.Fatal(e)
	}
	defer func() {
		if e = os.RemoveAll(dbPath); E.Chk(e) {
		}
	}()
	defer func() {
		if e = db.Close(); E.Chk(e) {
		}
	}()
	e = db.Update(func(tx database.Tx) (e error) {
		block := block2.NewBlock(chaincfg.MainNetParams.GenesisBlock)
		return tx.StoreBlock(block)
	},
	)
	if e != nil {
		b.Fatal(e)
	}
	b.ReportAllocs()
	b.ResetTimer()
	e = db.View(func(tx database.Tx) (e error) {
		blockHash := chaincfg.MainNetParams.GenesisHash
		for i := 0; i < b.N; i++ {
			_, e := tx.FetchBlockHeader(blockHash)
			if e != nil {
				return e
			}
		}
		return nil
	},
	)
	if e != nil {
		b.Fatal(e)
	}
	// Don't benchmark teardown.
	b.StopTimer()
}

// BenchmarkBlockHeader benchmarks how long it takes to load the mainnet genesis block.
func BenchmarkBlock(b *testing.B) {
	// Start by creating a new database and populating it with the mainnet genesis block.
	dbPath := filepath.Join(os.TempDir(), "ffldb-benchblk")
	_ = os.RemoveAll(dbPath)
	db, e := database.Create("ffldb", dbPath, blockDataNet)
	if e != nil {
		b.Fatal(e)
	}
	defer func() {
		if e = os.RemoveAll(dbPath); E.Chk(e) {
		}
	}()
	defer func() {
		if e = db.Close(); E.Chk(e) {
		}
	}()
	e = db.Update(func(tx database.Tx) (e error) {
		block := block2.NewBlock(chaincfg.MainNetParams.GenesisBlock)
		return tx.StoreBlock(block)
	},
	)
	if e != nil {
		b.Fatal(e)
	}
	b.ReportAllocs()
	b.ResetTimer()
	e = db.View(func(tx database.Tx) (e error) {
		blockHash := chaincfg.MainNetParams.GenesisHash
		for i := 0; i < b.N; i++ {
			_, e := tx.FetchBlock(blockHash)
			if e != nil {
				return e
			}
		}
		return nil
	},
	)
	if e != nil {
		b.Fatal(e)
	}
	// Don't benchmark teardown.
	b.StopTimer()
}
