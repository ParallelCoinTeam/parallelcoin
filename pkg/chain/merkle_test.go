package blockchain

import (
	"testing"

	"github.com/p9c/pod/pkg/util"
)

// TestMerkle tests the BuildMerkleTreeStore API.
func TestMerkle(t *testing.T) {
	block := util.NewBlock(&Block100000)
	merkles := BuildMerkleTreeStore(block.Transactions(), false)
	calculatedMerkleRoot := merkles.GetRoot()
	wantMerkle := &Block100000.Header.MerkleRoot
	if !wantMerkle.IsEqual(calculatedMerkleRoot) {
		t.Errorf("BuildMerkleTreeStore: merkle root mismatch - got %v, want %v", calculatedMerkleRoot, wantMerkle)
	}
}
