package spv

import (
	"fmt"
	"github.com/stalker-loki/app/slog"

	"github.com/p9c/pod/cmd/spv/headerfs"
	blockchain "github.com/p9c/pod/pkg/chain"
	chainhash "github.com/p9c/pod/pkg/chain/hash"
	"github.com/p9c/pod/pkg/chain/wire"
	waddrmgr "github.com/p9c/pod/pkg/wallet/addrmgr"
)

// mockBlockHeaderStore is an implementation of the BlockHeaderStore backed by
// a simple map.
type mockBlockHeaderStore struct {
	headers map[chainhash.Hash]wire.BlockHeader
}

// A compile-time check to ensure the mockBlockHeaderStore adheres to the
// BlockHeaderStore interface.
var _ headerfs.BlockHeaderStore = (*mockBlockHeaderStore)(nil)

// NewMockBlockHeaderStore returns a version of the BlockHeaderStore that's
// backed by an in-memory map. This instance is meant to be used by callers
// outside the package to unit test components that require a BlockHeaderStore
// interface.
// nolint
func newMockBlockHeaderStore() headerfs.BlockHeaderStore {
	return &mockBlockHeaderStore{
		headers: make(map[chainhash.Hash]wire.BlockHeader),
	}
}
func (m *mockBlockHeaderStore) ChainTip() (b *wire.BlockHeader, u uint32, err error) {
	return nil, 0, nil
}
func (m *mockBlockHeaderStore) LatestBlockLocator() (loc blockchain.BlockLocator, err error) {
	return nil, nil
}
func (m *mockBlockHeaderStore) FetchHeaderByHeight(height uint32) (b *wire.BlockHeader, err error) {
	return nil, nil
}
func (m *mockBlockHeaderStore) FetchHeaderAncestors(uint32, *chainhash.Hash) (b []wire.BlockHeader, u uint32, err error) {
	return nil, 0, nil
}
func (m *mockBlockHeaderStore) HeightFromHash(*chainhash.Hash) (u uint32, err error) {
	return 0, nil
}
func (m *mockBlockHeaderStore) RollbackLastBlock() (bs *waddrmgr.BlockStamp, err error) {
	return nil, nil
}
func (m *mockBlockHeaderStore) FetchHeader(h *chainhash.Hash) (header *wire.BlockHeader, u uint32, err error) {
	var (
		ok  bool
		hdr wire.BlockHeader
	)

	if hdr, ok = m.headers[*h]; ok {
		header = &hdr
		return
	}
	err = fmt.Errorf("not found")
	slog.Debug(err)
	return
}
func (m *mockBlockHeaderStore) WriteHeaders(headers ...headerfs.BlockHeader) (err error) {
	for _, h := range headers {
		m.headers[h.BlockHash()] = *h.BlockHeader
	}
	return
}
