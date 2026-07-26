# 01 — Worker Pool with Backpressure + Cancellation

## Context

For every catering search we need to score ~10k restaurants in <200ms.
Each restaurant has lat/lon, rating, price. The user sends their lat/lon.
The score combines distance + rating + relevance. Return top 20.

## Signatures (you define the types)

```
type Config struct { … }
type Ranker struct { … }
func New(cfg Config) *Ranker
func (r *Ranker) Rank(ctx context.Context, restaurants []Restaurant, userLat, userLon float64) ([]ScoredRestaurant, error)
```

## Hints

- Limit concurrency with `golang.org/x/sync/errgroup` + `semaphore.Weighted` (or a `chan struct{}` semaphore)
- `errgroup` cancels ctx on first error — check `ctx.Err()` at the start of each worker
- Top-N: use `container/heap` (max-heap of size TopN) — O(n log k), not O(n log n)
- Score function: haversine + rating weight + relevance weight. Add a tiny `time.Sleep` to simulate work
- Prepare `BenchmarkRanker10k` — generate 10k restaurants, measure time + allocs

## Tests

- `TestRanker_ReturnsExactN` — 20 in, TopN=3 → 3 results, sorted desc
- `TestRanker_CancelledContext` — ctx already cancelled → immediate ctx.Err()
- `TestRanker_Timeout` — deadline too short → deadline exceeded
- `TestRanker_Backpressure` — MaxWorkers=1, 50 restaurants → verify ≤1 goroutine
- `BenchmarkRanker10k` — 10k restaurants, TopN=20

## Alternatives

| Approach | Pros | Cons |
|----------|------|------|
| errgroup + semaphore | Simple, native cancellation, readable | Manual backpressure |
| sync.WaitGroup | Full lifecycle control | Manual cancellation, more boilerplate |
| Channel pipeline (gen → workers → merger → topN) | Composable, testable per stage | Channel overhead, more complex |
| singleflight.Group | Collapses duplicate requests | Complementary (caching layer), not a replacement |

## Interview angle

- How do you handle stragglers?
- Horizontal scaling? (shard by region, precompute base scores)
- Where does caching fit? (data layer, not worker)
- If the scoring model is an ML service? (latency, circuit breaker, retry)
- Why heap and not sort.Slice?
