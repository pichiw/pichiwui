package entity

import (
	"sync"
	"time"

	"github.com/gopherjs/gopherwasm/js"
	"github.com/gowasm/vecty"
	"github.com/pichiw/leaflet"
	"github.com/pichiw/pichiwui/model"
)

// OnEntityClick is called when an entity is clicked
type OnEntityClick func(e *model.Entity)

func NewMap(
	m *leaflet.Map,
	onClick OnEntityClick,
	entities ...*model.Entity,
) *Map {
	var markers []*leaflet.Marker

	for _, e := range entities {
		markers = append(markers, leaflet.NewMarker(e.Coord, leaflet.Events{"click": onClicker(onClick, e)}))
	}

	valuers := make([]vecty.JSValuer, 0, len(markers))
	for _, m := range markers {
		valuers = append(valuers, m)
	}
	return &Map{
		m:        m,
		entities: entities,
		markers:  markers,
		group:    leaflet.NewFeatureGroup(valuers...),
	}
}

type Map struct {
	m         *leaflet.Map
	entities  []*model.Entity
	markers   []*leaflet.Marker
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

func (em *Map) Show(start, end time.Time) {
	em.showMutex.Lock()
	defer em.showMutex.Unlock()

	initial := false
	if em.shown == nil {
		em.shown = make([]bool, len(em.entities))
		initial = true
	}

	for i, e := range em.entities {
		shown := !(e.Time.Before(start) || e.Time.After(end))

		currentShown := em.shown[i]
		if shown == currentShown {
			continue
		}

		em.shown[i] = shown

		if shown {
			em.markers[i].AddTo(em.m)
		} else if !initial {
			em.markers[i].Remove()
		}
	}
}
