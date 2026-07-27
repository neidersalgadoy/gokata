package kata

import (
	"context"
	"sync/atomic"
	"testing"
	"time"
)

func TestCache_Hit(t *testing.T) {
	var calls atomic.Int32
	svc := &fakeScorer{
		fn: func(ctx context.Context, lat, lon float64) ([]ScoredRestaurant, error) {
			calls.Add(1)
			return result(lat, lon), nil
		},
	}

	c := NewCache(svc, 10*time.Minute)
	_, err := c.Get(context.Background(), 40.71, -74.00)
	if err != nil {
		t.Fatal(err)
	}
	_, err = c.Get(context.Background(), 40.71, -74.00)
	if err != nil {
		t.Fatal(err)
	}
	if calls.Load() != 1 {
		t.Fatalf("expected 1 backend call, got %d", calls.Load())
	}
}

func TestCache_Miss(t *testing.T) {
	var calls atomic.Int32
	svc := &fakeScorer{
		fn: func(ctx context.Context, lat, lon float64) ([]ScoredRestaurant, error) {
			calls.Add(1)
			return result(lat, lon), nil
		},
	}

	c := NewCache(svc, 10*time.Minute)
	_, err := c.Get(context.Background(), 40.71, -74.00)
	if err != nil {
		t.Fatal(err)
	}
	if calls.Load() != 1 {
		t.Fatalf("expected 1 backend call, got %d", calls.Load())
	}
}

func TestCache_Dedup(t *testing.T) {
	var calls atomic.Int32
	svc := &fakeScorer{
		fn: func(ctx context.Context, lat, lon float64) ([]ScoredRestaurant, error) {
			calls.Add(1)
			time.Sleep(50 * time.Millisecond)
			return result(lat, lon), nil
		},
	}

	c := NewCache(svc, 10*time.Minute)

	ctx := context.Background()
	const n = 10
	errs := make(chan error, n)
	for range n {
		go func() {
			_, err := c.Get(ctx, 40.71, -74.00)
			errs <- err
		}()
	}

	for range n {
		if err := <-errs; err != nil {
			t.Fatal(err)
		}
	}

	if calls.Load() != 1 {
		t.Fatalf("expected 1 backend call for 10 concurrent requests, got %d", calls.Load())
	}
}

func TestCache_DifferentKeys(t *testing.T) {
	var calls atomic.Int32
	svc := &fakeScorer{
		fn: func(ctx context.Context, lat, lon float64) ([]ScoredRestaurant, error) {
			calls.Add(1)
			return result(lat, lon), nil
		},
	}

	c := NewCache(svc, 10*time.Minute)

	locs := []struct{ lat, lon float64 }{
		{40.71, -74.00},
		{34.05, -118.24},
		{51.50, -0.12},
	}

	for _, loc := range locs {
		_, err := c.Get(context.Background(), loc.lat, loc.lon)
		if err != nil {
			t.Fatal(err)
		}
	}

	if calls.Load() != 3 {
		t.Fatalf("expected 3 backend calls for 3 different keys, got %d", calls.Load())
	}
}

func TestCache_Expiry(t *testing.T) {
	var calls atomic.Int32
	svc := &fakeScorer{
		fn: func(ctx context.Context, lat, lon float64) ([]ScoredRestaurant, error) {
			calls.Add(1)
			return result(lat, lon), nil
		},
	}

	c := NewCache(svc, 50*time.Millisecond)

	_, err := c.Get(context.Background(), 40.71, -74.00)
	if err != nil {
		t.Fatal(err)
	}

	time.Sleep(100 * time.Millisecond)

	_, err = c.Get(context.Background(), 40.71, -74.00)
	if err != nil {
		t.Fatal(err)
	}

	if calls.Load() != 2 {
		t.Fatalf("expected 2 backend calls (miss + expiry refetch), got %d", calls.Load())
	}
}

func BenchmarkCacheContention(b *testing.B) {
	svc := &fakeScorer{
		fn: func(ctx context.Context, lat, lon float64) ([]ScoredRestaurant, error) {
			time.Sleep(10 * time.Millisecond)
			return result(lat, lon), nil
		},
	}

	c := NewCache(svc, 10*time.Minute)
	ctx := context.Background()

	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			c.Get(ctx, 40.71, -74.00)
		}
	})
}

// --- fakes ---

type fakeScorer struct {
	fn func(ctx context.Context, lat, lon float64) ([]ScoredRestaurant, error)
}

func (f *fakeScorer) Score(ctx context.Context, lat, lon float64) ([]ScoredRestaurant, error) {
	return f.fn(ctx, lat, lon)
}

func result(lat, lon float64) []ScoredRestaurant {
	return []ScoredRestaurant{
		{Score: 95, Restaurant: Restaurant{ID: 1, Lat: lat, Lon: lon, Rating: 4.5}},
		{Score: 85, Restaurant: Restaurant{ID: 2, Lat: lat, Lon: lon, Rating: 4.0}},
	}
}
