package kata

type Restaurant struct {
	ID        int
	Lat       float64
	Lon       float64
	Rating    float64
	Relevance float64
}

type ScoredRestaurant struct {
	Restaurant
	Score float64
}

type Config struct {
	MaxWorkers int
	TopN       int
}
