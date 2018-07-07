package main

import (
	"fmt"
	"image/color"
	"sync"

	"github.com/gowasm/vecty"
	"github.com/gowasm/vecty/elem"
	"github.com/pichiw/leaflet"
	"github.com/pichiw/md"
	"github.com/pichiw/pichiwui/components"
)

const mapDiv = "mapid"

func main() {
	c := make(chan struct{}, 0)

	model := []*components.Entity{
		&components.Entity{
			Name:  "Red River",
			Coord: leaflet.NewCoordinate(49.8951, -97.1384),
		},
		&components.Entity{
			Name:  "Turtle Mountain",
			Coord: leaflet.NewCoordinate(48.8469, -99.8011),
		},
		&components.Entity{
			Name:  "St. Paul des Metis",
			Coord: leaflet.NewCoordinate(53.8896, -111.4657),
		},
		&components.Entity{
			Name:  "Vancouver",
			Coord: leaflet.NewCoordinate(49.2827, -123.1207),
		},
	}

	b := &Home{
		model:  model,
		mapDiv: mapDiv,
		min:    0,
		max:    len(model) - 1,
	}

	vecty.SetTitle("Pichiw")
	vecty.AddStylesheet("app.css")
	vecty.RenderBody(b)

	<-c
}

type Home struct {
	vecty.Core

	mapDiv    string
	m         *leaflet.Map
	entityMap *components.EntityMap
	mapOnce   sync.Once

	currentEntity *components.Entity
	entityMutex   sync.Mutex

	slider *md.Slider
	model  []*components.Entity

	min      int
	max      int
	valMutex sync.Mutex
}

func shown(model []*components.Entity, min, max int) []bool {
	s := make([]bool, len(model))
	for i := min; i < max; i++ {
		s[i] = true
	}
	return s
}

func (p *Home) onSliderChange(s *md.Slider) {
	p.valMutex.Lock()
	p.max = int(s.Value())
	p.entityMap.Show(shown(p.model, p.min, p.max))
	p.valMutex.Unlock()
	vecty.Rerender(p)
}

func (p *Home) Mount() {
	p.mapOnce.Do(func() {
		p.m = leaflet.NewMap(
			p.mapDiv,
			leaflet.MapOptions{
				Center:  leaflet.NewCoordinate(49.8951, -97.1384),
				Zoom:    6,
				MaxZoom: 18,
			},
			nil,
		)
		p.m.Add(
			leaflet.NewTileLayer(
				leaflet.TileLayerOptions{
					MaxZoom:     18,
					Attribution: `&copy; <a href="http://www.openstreetmap.org/copyright">OpenStreetMap</a>`,
				},
			),
		)

		p.entityMap = components.NewEntityMap(p.m, p.onEntityClick, color.RGBA{255, 0, 0, 255}, p.model...)
		p.entityMap.Show(shown(p.model, 0, len(p.model)))
	})
}

func (p *Home) onEntityClick(e *components.Entity) {
	p.entityMutex.Lock()
	p.currentEntity = e
	p.entityMutex.Unlock()
	vecty.Rerender(p)
	fmt.Println(e.Name)
}

// Render implements the vecty.Component interface.
func (p *Home) Render() vecty.ComponentOrHTML {
	p.entityMutex.Lock()
	hasElement := p.currentEntity != nil
	p.entityMutex.Unlock()
	mapSpan := 10
	if hasElement {
		mapSpan = 6
	}
	return elem.Body(
		md.LayoutGrid(
			md.LayoutGridInner(
				md.LayoutGridCell(md.LayoutGridCellOptions{Span: 2},
					elem.Heading1(vecty.Text("Pichiw")),
				),
				md.LayoutGridCell(md.LayoutGridCellOptions{Span: mapSpan},
					elem.Div(vecty.Markup(vecty.Attribute("id", p.mapDiv))),
				),
				vecty.If(hasElement,
					elem.Div(
						md.LayoutGridCell(
							md.LayoutGridCellOptions{Span: 4},
							components.NewEntityEditor(p.currentEntity),
						),
					),
				),
				md.LayoutGridCell(md.LayoutGridCellOptions{Span: 12},
					md.NewSlider(
						md.SliderOptions{
							Min:      0,
							Max:      float64(len(p.model)),
							OnChange: p.onSliderChange,
							OnInput:  p.onSliderChange,
						},
					),
				),
			),
		),
	)
}
