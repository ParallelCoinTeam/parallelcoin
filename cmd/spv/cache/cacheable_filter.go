package cache

import (
	"github.com/p9c/pod/pkg/coding/gcs"
	"github.com/stalker-loki/app/slog"
)

// CacheableFilter is a wrapper around Filter type which provides a Size method
// used by the cache to target certain memory usage.
type CacheableFilter struct {
	*gcs.Filter
}

// Size returns size of this filter in bytes.
func (c *CacheableFilter) Size() (uint64, error) {
	f, err := c.Filter.NBytes()
	if err != nil {
		slog.Error(err)
		return 0, err
	}
	return uint64(len(f)), nil
}
