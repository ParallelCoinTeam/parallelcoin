package cache

import (
	"github.com/p9c/pod/pkg/block"
)

// CacheableBlock is a wrapper around the util.Block type which provides a Size method used by the cache to target
// certain memory usage.
type CacheableBlock struct {
	*block.Block
}

// Size returns size of this block in bytes.
func (c *CacheableBlock) Size() (rv uint64, e error) {
	return uint64(c.Block.WireBlock().SerializeSize()), nil
}
