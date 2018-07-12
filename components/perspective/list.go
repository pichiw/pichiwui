package perspective

import (
	"sort"
	"sync"
	"time"

	"github.com/gowasm/vecty"
	"github.com/gowasm/vecty/elem"
	"github.com/pichiw/leaflet"
	"github.com/pichiw/md"
	"github.com/pichiw/pichiwui/components/entity"
	"github.com/pichiw/pichiwui/model"
)

func NewList(perspective *model.Perspective) *List {
	entities := perspective.AllEntities()
	sort.Sort(model.EntitySort(entities))
	return &List{
		perspective: perspective,
		entities:    entities,
		mapDiv:      "mapid",
		min:         entities[0].Time,
		currentMax:  entities[len(entities)-1].Time,
		max:         entities[len(entities)-1].Time,
		editor:      entity.NewEditor(),
	}
}

type List struct {
	vecty.Core

	mapDiv  string
	m       *leaflet.Map
	em      *entity.Map
	pm      *Map
	mapOnce sync.Once

	editor *entity.Editor

	slider      *md.Slider
	entities    []*model.Entity
	perspective *model.Perspective

	min        time.Time
	currentMax time.Time
	max        time.Time
	valMutex   sync.Mutex
}

func (l *List) onSliderChange(s *md.Slider) {
	l.valMutex.Lock()
	newValue := l.min.Add(time.Duration(s.Value()) * 7 * 24 * time.Hour)
	if newValue == l.currentMax {
		l.valMutex.Unlock()
		return
	}
	l.currentMax = newValue
	l.pm.Show(l.min, l.currentMax)
	l.em.Show(l.min, l.currentMax)
	l.valMutex.Unlock()
	vecty.Rerender(l)
}

func (l *List) Mount() {
	l.mapOnce.Do(func() {
		l.m = leaflet.NewMap(
			l.mapDiv,
			leaflet.MapOptions{
				Center:  leaflet.NewCoordinate(49.8951, -97.1384),
				Zoom:    6,
				MaxZoom: 18,
			},
			nil,
		)

		leaflet.NewTileLayer(
			leaflet.TileLayerOptions{
				MaxZoom:     18,
				Attribution: `&copy; <a href="http://www.openstreetmal.org/copyright">OpenStreetMap</a>`,
			},
		).AddTo(l.m)

		l.pm = NewMap(l.m, l.perspective)
		l.pm.Show(l.min, l.currentMax)
		l.em = entity.NewMap(l.m, l.onEntityClick, l.entities...)
		l.em.Show(l.min, l.currentMax)
		l.em.Bounds()
	})
}

func (l *List) onEntityClick(e *model.Entity) {
	if l.editor.Entity() == e {
		l.editor.SetEntity(nil)
	} else {
		l.editor.SetEntity(e)
	}
	vecty.Rerender(l)
}

// Render implements the vecty.Component interface.
func (l *List) Render() vecty.ComponentOrHTML {
	hasElement := l.editor.Entity() != nil

	mapSpan := 12
	if hasElement {
		mapSpan = 8
	}

	delta := l.max.Sub(l.min).Hours()
	delta /= 24 * 7

	return elem.Body(
		md.LayoutGrid(
			md.LayoutGridInner(
				md.LayoutGridCell(md.LayoutGridCellOptions{Span: 12},
					elem.Heading1(vecty.Text("Pichiw")),
				),
				md.LayoutGridCell(md.LayoutGridCellOptions{Span: mapSpan},
					md.LayoutGrid(
						md.LayoutGridInner(
							md.LayoutGridCell(md.LayoutGridCellOptions{Span: 12},
								elem.Div(vecty.Markup(vecty.Attribute("id", l.mapDiv))),
							),
							md.LayoutGridCell(md.LayoutGridCellOptions{Span: 1},
								elem.Heading2(vecty.Text(l.min.Format("2006-01-02"))),
							),
							md.LayoutGridCell(md.LayoutGridCellOptions{Span: 10},
								md.LayoutGrid(
									md.LayoutGridInner(
										md.LayoutGridCell(md.LayoutGridCellOptions{Span: 12},
											md.NewSlider(
												md.SliderOptions{
													Min:      0,
													Max:      float64(delta),
													OnChange: l.onSliderChange,
													OnInput:  l.onSliderChange,
												},
											),
										),
										md.LayoutGridCell(md.LayoutGridCellOptions{Span: 12},
											elem.Heading2(vecty.Text("Selected: "+l.currentMax.Format("2006-01-02"))),
										),
									),
								),
							),
							md.LayoutGridCell(md.LayoutGridCellOptions{Span: 1},
								elem.Heading2(vecty.Text(l.max.Format("2006-01-02"))),
							),
						),
					),
				),
				vecty.If(hasElement,
					elem.Div(
						md.LayoutGridCell(
							md.LayoutGridCellOptions{Span: 4},
							l.editor,
						),
					),
				),
			),
		),
	)
}
