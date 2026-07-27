# 02 — Cache Stampede Prevention with Singleflight

## Context

The restaurant ranking service caches top-20 results per (lat,lon) with a 30s TTL.
When the cache expires and 100 search requests arrive simultaneously for the same location,
all 100 hit the backend scoring service. This cache stampede melts the database.

Singleflight collapses concurrent requests for the same key into one backend call.
One request does the work; the rest share the result.

## Signatures

```
type Scorer interface {
    Score(ctx context.Context, userLat, userLon float64) ([]ScoredRestaurant, error)
}

type Cache struct { ... }
func NewCache(svc Scorer, ttl time.Duration) *Cache
func (c *Cache) Get(ctx context.Context, userLat, userLon float64) ([]ScoredRestaurant, error)
```

## Hints

- Use `golang.org/x/sync/singleflight.Group` — the `Do` method is the star
- Key = `fmt.Sprintf("%.4f:%.4f", lat, lon)` — 4 decimals to control cache granularity
- In-memory cache: `sync.Map` or `map[string]cacheEntry` + `sync.RWMutex`
- TTL: store `expiresAt` in each entry; check before returning
- On cache hit → return immediately, no singleflight call
- On cache miss → `g.Do(key, fn)` where `fn` calls `svc.Score` and stores in cache
- Singleflight caveat: `Do` returns to ALL callers — they all share the same result
- Context: singleflight doesn't natively support per-call context cancellation. Use `DoChan` for that.

## Tests

- `TestCache_Hit` — cached value returns without calling Scorer
- `TestCache_Miss` — first call goes to Scorer, result cached
- `TestCache_Dedup` — 10 concurrent requests for same (lat,lon) → Scorer called once
- `TestCache_DifferentKeys` — 2 different locations → 2 separate Scorer calls
- `TestCache_Expiry` — stale entry is re-fetched
- `BenchmarkCacheContention` — 100 concurrent requests, measure dedup ratio

## Alternatives

| Approach | Pros | Cons |
|----------|------|------|
| singleflight.Group | Simple, proven, native dedup | No per-call ctx cancel (use DoChan) |
| Mutex + dedup map + WaitGroup | Full control, no dependency | Boilerplate, race-prone |
| Redis locking (SETNX) | Distributed, across instances | Latency, complexity, not needed for single-node |
| Probabilistic early expiry (XFetch) | Proactive, no stampede | Tuning needed, probabilistic |

## Interview angle

- What happens if singleflight's shared result is stale? (TTL + background refresh)
- How do you handle thundering herd across N servers? (distributed singleflight via Redis)
- What if the backend call returns an error? (singleflight retries, cache last-good-value)
- How does singleflight compare to a mutex-based dedup?
- What's the memory behavior of singleflight? (Group holds all in-flight calls)
