package entity

import (
	"image/color"
	"sync"

	"github.com/gowasm/vecty"
	"github.com/gowasm/vecty/elem"
	"github.com/pichiw/leaflet"
	"github.com/pichiw/md"
	"github.com/pichiw/pichiwui/model"
)

func NewList(e []*model.Entity) *List {
	return &List{
		model:  e,
		mapDiv: "mapid",
		min:    0,
		max:    len(e),
		editor: NewEditor(),
	}
}

type List struct {
	vecty.Core

	mapDiv  string
	m       *leaflet.Map
	em      *Map
	mapOnce sync.Once

	editor *Editor

	slider *md.Slider
	model  []*model.Entity

	min      int
	max      int
	valMutex sync.Mutex
}

func shown(model []*model.Entity, min, max int) []bool {
	s := make([]bool, len(model))
	for i := min; i < max; i++ {
		s[i] = true
	}
	return s
}

func (l *List) onSliderChange(s *md.Slider) {
	l.valMutex.Lock()
	newValue := int(s.Value())
	if newValue == l.max {
		l.valMutex.Unlock()
		return
	}
	l.max = newValue
	l.em.Show(shown(l.model, l.min, l.max))
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

		l.em = NewMap(l.m, l.onEntityClick, color.RGBA{255, 0, 0, 255}, l.model...)
		l.em.Show(shown(l.model, 0, len(l.model)))
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
							md.LayoutGridCell(md.LayoutGridCellOptions{Span: 12},
								md.NewSlider(
									md.SliderOptions{
										Min:      float64(l.min),
										Max:      float64(l.max),
										OnChange: l.onSliderChange,
										OnInput:  l.onSliderChange,
									},
								),
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
