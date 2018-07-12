package model

import (
	"time"

	"github.com/pichiw/leaflet"
)

type Entity struct {
	Name  string
	Time  time.Time
	Coord *leaflet.Coordinate
}

// EntitySort is a type for sorting entities
type EntitySort []*Entity

// Len is part of sort.Interface.
func (s EntitySort) Len() int {
	return len(s)
}

// Swap is part of sort.Interface.
func (s EntitySort) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}

// Less is part of sort.Interface. It is implemented by calling the "by" closure in the sorter.
func (s EntitySort) Less(i, j int) bool {
	return s[i].Time.Before(s[j].Time)
}
