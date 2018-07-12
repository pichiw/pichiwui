package perspective

import (
	"sync"
	"time"

	"github.com/pichiw/leaflet"
	"github.com/pichiw/pichiwui/htmlhelp"
	"github.com/pichiw/pichiwui/model"
)

// OnEntityClick is called when an entity is clicked
type OnEntityClick func(e *model.Entity)

func NewMap(
	m *leaflet.Map,
	perspective *model.Perspective,
) *Map {
	var cw htmlhelp.ColorWheel
	return &Map{
		m:          m,
		polylines:  perspectivePolylines(perspective, nil, &cw),
		validTimes: perspectiveValidTimes(perspective, nil),
	}
}

func perspectiveValidTimes(perspective *model.Perspective, last *model.Entity) []time.Time {
	var vts []time.Time

	for _, e := range perspective.Entities {
		if last == nil {
			last = e
			continue
		}
		vts = append(vts, e.Time)
		last = e
	}

	for _, c := range perspective.Children {
		vts = append(vts, perspectiveValidTimes(c, last)...)
	}

	return vts
}

func perspectivePolylines(perspective *model.Perspective, last *model.Entity, cw *htmlhelp.ColorWheel) []*leaflet.Polyline {
	var polylines []*leaflet.Polyline

	color := cw.NextColor()
	color.A = 128

	for _, e := range perspective.Entities {
		if last == nil {
			last = e
			continue
		}
		polylines = append(polylines,
			leaflet.NewPolyline(
				leaflet.PolylineOptions{
					PathOptions: leaflet.PathOptions{
						Color: htmlhelp.HTMLColor(color),
					},
				},
				last.Coord,
				e.Coord,
			),
		)
		last = e
	}

	for _, c := range perspective.Children {
		polylines = append(polylines, perspectivePolylines(c, last, cw)...)
	}

	return polylines
}

type Map struct {
	m           *leaflet.Map
	perspective *model.Perspective
	polylines   []*leaflet.Polyline
	validTimes  []time.Time

	showMutex sync.Mutex
	shown     []bool
}

func (pm *Map) Show(start, end time.Time) {
	pm.showMutex.Lock()
	defer pm.showMutex.Unlock()

	initial := false
	if pm.shown == nil {
		pm.shown = make([]bool, len(pm.polylines))
		initial = true
	}

	for i, vt := range pm.validTimes {
		shown := !(vt.Before(start) || vt.After(end))

		currentShown := pm.shown[i]
		if shown == currentShown {
			continue
		}

		pm.shown[i] = shown

		if shown {
			pm.polylines[i].AddTo(pm.m)
		} else if !initial {
			pm.polylines[i].Remove()
		}
	}
}
