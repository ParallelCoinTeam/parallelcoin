package blockchain

import (
	"fmt"
	"runtime"
	"testing"
	
	"github.com/p9c/matrjoska/pkg/txscript"
)

// TestCheckBlockScripts ensures that validating the all of the scripts in a known-good block doesn't return an error.
func TestCheckBlockScripts(t *testing.T) {
	runtime.GOMAXPROCS(runtime.NumCPU())
	testBlockNum := 277647
	blockDataFile := fmt.Sprintf("%d.dat.bz2", testBlockNum)
	blocks, e := loadBlocks(blockDataFile)
	if e != nil {
		t.Errorf("ScriptError loading file: %v\n", e)
		return
	}
	if len(blocks) > 1 {
		t.Errorf("The test block file must only have one block in it")
		return
	}
	if len(blocks) == 0 {
		t.Errorf("The test block file may not be empty")
		return
	}
	storeDataFile := fmt.Sprintf("%d.utxostore.bz2", testBlockNum)
	view, e := loadUtxoView(storeDataFile)
	if e != nil {
		t.Errorf("ScriptError loading txstore: %v\n", e)
		return
	}
	scriptFlags := txscript.ScriptBip16
	e = checkBlockScripts(blocks[0], view, scriptFlags, nil, nil)
	if e != nil {
		t.Errorf("Transaction script validation failed: %v\n", e)
		return
	}
}
