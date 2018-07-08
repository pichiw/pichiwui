package entity

import (
	"image/color"
	"sync"

	"github.com/gopherjs/gopherwasm/js"
	"github.com/pichiw/leaflet"
	"github.com/pichiw/pichiwui/components"
)

// OnEntityClick is called when an entity is clicked
type OnEntityClick func(e *Entity)

func NewMap(
	m *leaflet.Map,
	onClick OnEntityClick,
	lineColor color.RGBA,
	entities ...*Entity,
) *Map {
	var markers []*leaflet.Marker
	var polylines []*leaflet.Polyline

	for i, e := range entities {
		markers = append(markers, leaflet.NewMarker(e.Coord, leaflet.Events{"click": onClicker(onClick, e)}))

		if i > 0 {
			polylines = append(polylines,
				leaflet.NewPolyline(
					leaflet.PolylineOptions{
						PathOptions: leaflet.PathOptions{
							Color: components.HTMLColor(lineColor),
						},
					},
					entities[i-1].Coord,
					e.Coord,
				),
			)
		}
	}

	return &Map{
		m:         m,
		entities:  entities,
		markers:   markers,
		polylines: polylines,
	}
}

type Map struct {
	m         *leaflet.Map
	entities  []*Entity
	markers   []*leaflet.Marker
	polylines []*leaflet.Polyline
	callbacks []js.Callback
	onClick   OnEntityClick
	start     int
	end       int

	showMutex sync.Mutex
	shown     []bool
}

func onClicker(onClick OnEntityClick, e *Entity) func(vs []js.Value) {
	return func(vs []js.Value) { onClick(e) }
}

func (em *Map) Show(shown []bool) {
	if len(shown) != len(em.entities) {
		panic("invalid shown length")
	}

	em.showMutex.Lock()
	defer em.showMutex.Unlock()

	initial := false
	if em.shown == nil {
		em.shown = make([]bool, len(shown))
		initial = true
	}

	for i, s := range shown {
		if s == em.shown[i] {
			continue
		}

		if s {
			em.m.Add(em.markers[i])
			if i > 0 {
				em.m.Add(em.polylines[i-1])
			}
		} else {
			if !initial {
				em.markers[i].Remove()
				if i > 0 {
					em.polylines[i-1].Remove()
				}
			}
		}
	}

	em.shown = shown
}
