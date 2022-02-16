package main

import "math/rand"

// SomeData represent some real data of some sort, unaware of tables
type SomeData struct {
	ID     string
	Score  int
	Status string
}

// NewSomeData creates SomeData that has an ID and randomized values
func NewSomeData(id string) *SomeData {
	s := &SomeData{
		ID: id,
	}

	// Start with some random data
	s.RandomizeScoreAndStatus()

	return s
}

// RandomizeScoreAndStatus does an in-place update to simulate some data being
// updated by some other process
func (s *SomeData) RandomizeScoreAndStatus() {
	s.Score = rand.Intn(100) + 1

	if s.Score < 30 {
		s.Status = "Critical"
	} else if s.Score < 80 {
		s.Status = "Stable"
	} else {
		s.Status = "Good"
	}
}
