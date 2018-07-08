package entity

import "github.com/pichiw/leaflet"

type Entity struct {
	Name  string
	Coord *leaflet.Coordinate
}
