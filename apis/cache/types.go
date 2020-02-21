package cache

// CacheFlushResult represent the result of a cache-flush.
type CacheFlushResult struct {
	Count  int    `json:"count"`
	Result string `json:"result"`
}
