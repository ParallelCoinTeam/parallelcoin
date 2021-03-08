package cache

import (
	"github.com/p9c/pod/pkg/coding/gcs"
)

// CacheableFilter is a wrapper around Filter type which provides a Size method used by the cache to target certain
// memory usage.
type CacheableFilter struct {
	*gcs.Filter
}

// Size returns size of this filter in bytes.
func (c *CacheableFilter) Size() (rv uint64, e error) {
	var f []byte
	f, e = c.Filter.NBytes()
	if e != nil {
		return 0, e
	}
	return uint64(len(f)), nil
}
