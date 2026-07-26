package kata

import "context"

type Ranker struct {
	cfg Config
}

func New(cfg Config) *Ranker {
	panic("implement me")
}

func (r *Ranker) Rank(ctx context.Context, restaurants []Restaurant, userLat, userLon float64) ([]ScoredRestaurant, error) {
	panic("implement me")
}
