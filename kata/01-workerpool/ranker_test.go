package kata

import (
	"context"
	"math"
	"testing"
	"time"
)

func TestRanker_ReturnsExactN(t *testing.T) {
	restaurants := make([]Restaurant, 20)
	for i := range restaurants {
		restaurants[i] = Restaurant{
			ID:     i,
			Lat:    40.7128 + float64(i)*0.001,
			Lon:    -74.0060 + float64(i)*0.001,
			Rating: float64(i) / 20,
		}
	}

	r := New(Config{MaxWorkers: 4, TopN: 3})
	results, err := r.Rank(context.Background(), restaurants, 40.7128, -74.0060)
	if err != nil {
		t.Fatal(err)
	}
	if len(results) != 3 {
		t.Fatalf("want 3 results, got %d", len(results))
	}
	for i := 1; i < len(results); i++ {
		if results[i].Score > results[i-1].Score {
			t.Fatal("not sorted descending")
		}
	}
}

func TestRanker_CancelledContext(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	r := New(Config{MaxWorkers: 2, TopN: 5})
	_, err := r.Rank(ctx, []Restaurant{{ID: 1}}, 0, 0)
	if err == nil {
		t.Fatal("expected error on cancelled context")
	}
}

func TestRanker_Timeout(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Nanosecond)
	defer cancel()
	time.Sleep(time.Millisecond) // let deadline pass

	restaurants := make([]Restaurant, 100)
	for i := range restaurants {
		restaurants[i] = Restaurant{ID: i, Lat: 40.71, Lon: -74.00}
	}

	r := New(Config{MaxWorkers: 10, TopN: 5})
	_, err := r.Rank(ctx, restaurants, 40.71, -74.00)
	if err == nil {
		t.Fatal("expected deadline exceeded error")
	}
}

func TestRanker_Backpressure(t *testing.T) {
	restaurants := make([]Restaurant, 50)
	for i := range restaurants {
		restaurants[i] = Restaurant{ID: i, Lat: 40.71, Lon: -74.00}
	}

	r := New(Config{MaxWorkers: 1, TopN: 3})

	done := make(chan struct{})

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	go func() {
		results, err := r.Rank(ctx, restaurants, 40.71, -74.00)
		if err != nil {
			t.Log(err)
		}
		_ = results
		close(done)
	}()

	select {
	case <-done:
		// finished before we could check
	case <-time.After(1 * time.Second):
		t.Fatal("ranker appears to be deadlocked or too slow")
	}
}

func TestRanker_EmptyInput(t *testing.T) {
	r := New(Config{MaxWorkers: 2, TopN: 5})
	results, err := r.Rank(context.Background(), nil, 0, 0)
	if err != nil {
		t.Fatal(err)
	}
	if len(results) != 0 {
		t.Fatalf("want 0 results, got %d", len(results))
	}
}

func BenchmarkRanker10k(b *testing.B) {
	restaurants := make([]Restaurant, 10000)
	for i := range restaurants {
		restaurants[i] = Restaurant{
			ID:       i,
			Lat:      40.7128 + float64(i)*0.0001,
			Lon:      -74.0060 + float64(i)*0.0001,
			Rating:   float64(i) / 10000,
			Relevance: float64(i%100) / 100,
		}
	}

	r := New(Config{MaxWorkers: 4, TopN: 20})

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := r.Rank(context.Background(), restaurants, 40.7128, -74.0060)
		if err != nil {
			b.Fatal(err)
		}
	}
}

// --- helpers for test determinism ---

func haversineDeg(lat1, lon1, lat2, lon2 float64) float64 {
	dLat := (lat2 - lat1) * math.Pi / 180
	dLon := (lon2 - lon1) * math.Pi / 180
	a := math.Pow(math.Sin(dLat/2), 2) +
		math.Cos(lat1*math.Pi/180)*math.Cos(lat2*math.Pi/180)*
			math.Pow(math.Sin(dLon/2), 2)
	return 2 * 6371 * math.Asin(math.Sqrt(a))
}


