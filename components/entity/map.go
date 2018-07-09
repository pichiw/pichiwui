package entity

import (
	"image/color"
	"sync"

	"github.com/gopherjs/gopherwasm/js"
	"github.com/gowasm/vecty"
	"github.com/pichiw/leaflet"
	"github.com/pichiw/pichiwui/components"
	"github.com/pichiw/pichiwui/model"
)

// OnEntityClick is called when an entity is clicked
type OnEntityClick func(e *model.Entity)

func NewMap(
	m *leaflet.Map,
	onClick OnEntityClick,
	lineColor color.RGBA,
	entities ...*model.Entity,
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

	valuers := make([]vecty.JSValuer, 0, len(markers)+len(polylines))
	for _, m := range markers {
		valuers = append(valuers, m)
	}
	for _, p := range polylines {
		valuers = append(valuers, p)
	}
	return &Map{
		m:         m,
		entities:  entities,
		markers:   markers,
		polylines: polylines,
		group:     leaflet.NewFeatureGroup(valuers...),
	}
}

type Map struct {
	m         *leaflet.Map
	entities  []*model.Entity
	markers   []*leaflet.Marker
	polylines []*leaflet.Polyline
	callbacks []js.Callback
	group     *leaflet.FeatureGroup
	onClick   OnEntityClick
	start     int
	end       int

	showMutex sync.Mutex
	shown     []bool
}

func onClicker(onClick OnEntityClick, e *model.Entity) func(vs []js.Value) {
	return func(vs []js.Value) { onClick(e) }
}

func (em *Map) Bounds() {
	em.m.Bounds(em.group.Bounds())
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
			em.markers[i].AddTo(em.m)
			if i > 0 {
				em.polylines[i-1].AddTo(em.m)
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
