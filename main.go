package main

import (
	"sync"

	"github.com/gowasm/gopherwasm/js"

	"github.com/gowasm/vecty"
	"github.com/gowasm/vecty/elem"
	"github.com/pichiw/leaflet"
	"github.com/pichiw/md"
)

func main() {
	c := make(chan struct{}, 0)

	vecty.SetTitle("Pichiw")
	vecty.AddStylesheet("app.css")

	model := []*Entity{
		&Entity{
			Name:  "Red River",
			Coord: leaflet.NewCoordinate(49.8951, -97.1384),
		},
		&Entity{
			Name:  "Turtle Mountain",
			Coord: leaflet.NewCoordinate(48.8469, -99.8011),
		},
		&Entity{
			Name:  "St. Paul des Metis",
			Coord: leaflet.NewCoordinate(53.8896, -111.4657),
		},
		&Entity{
			Name:  "Vancouver",
			Coord: leaflet.NewCoordinate(49.2827, -123.1207),
		},
	}

	m := leaflet.NewMap(
		"mapid",
		leaflet.MapOptions{
			Center:  leaflet.NewCoordinate(49.8951, -97.1384),
			Zoom:    6,
			MaxZoom: 18,
		},
		nil,
	)
	b := &Home{
		model: model,
		m:     m,
		min:   0,
		max:   len(model) - 1,
	}

	onSliderChange := func(s *md.Slider) {
		b.valMutex.Lock()
		b.max = int(s.Value())
		b.valMutex.Unlock()
		b.updateMapFromModel()
		vecty.Rerender(b)
	}

	b.slider = md.NewSlider(
		md.SliderOptions{
			Min:      0,
			Max:      float64(len(model)),
			OnChange: onSliderChange,
			OnInput:  onSliderChange,
		},
	)

	vecty.RenderBody(b)

	<-c
}

type Entity struct {
	Name  string
	Coord *leaflet.Coordinate
}

type EntityEditor struct {
	vecty.Core
	entity *Entity
}

func (p *EntityEditor) Render() vecty.ComponentOrHTML {
	return elem.Div(
		elem.Heading2(
			vecty.Text(p.entity.Name),
		),
	)
}

type Home struct {
	vecty.Core
	m *leaflet.Map
	e EntityEditor

	slider  *md.Slider
	model   []*Entity
	polys   []*leaflet.Polyline
	markers []*leaflet.Marker

	min      int
	max      int
	valMutex sync.Mutex
}

func (p *Home) updateModel() {
	p.valMutex.Lock()
	defer p.valMutex.Unlock()

	for _, poly := range p.polys {
		poly.Remove()
	}

	for _, marker := range p.markers {
		marker.Remove()
	}

	p.polys = nil
	p.markers = nil

	colors := []string{
		"#ff0000",
		"#aa2200",
		"#66aa00",
		"#00ff00",
	}

	for i, entity := range p.model {
		e := entity
		if i > 0 {
			p.polys = append(p.polys,
				leaflet.NewPolyline(
					leaflet.PolylineOptions{
						PathOptions: leaflet.PathOptions{
							Color: colors[i],
						},
					},
					p.model[i-1].Coord,
					e.Coord,
				),
			)
		}

		p.markers = append(p.markers,
			leaflet.NewMarker(
				e.Coord,
				leaflet.Events{
					"click": func(vs []js.Value) {
						if p.e.entity == nil || p.e.entity != e {
							p.e.entity = e
							p.m.View(e.Coord, p.m.Zoom())
						} else {
							p.e.entity = nil
						}
						vecty.Rerender(p)
					},
				},
			),
		)
	}

	for _, m := range p.markers {
		p.m.Add(m)
	}

	for _, poly := range p.polys {
		p.m.Add(poly)
	}
}

func (p *Home) updateMapFromModel() {
	p.valMutex.Lock()
	defer p.valMutex.Unlock()

	for i := range p.model {
		if i > 0 {
			p.polys[i-1].Remove()
		}

		p.markers[i].Remove()
	}

	for i := range p.model {
		if i < p.min || i > p.max {
			continue
		}

		if i > 0 {
			p.m.Add(p.polys[i-1])
		}

		p.m.Add(p.markers[i])
	}
}

func (p *Home) Mount() {
	p.m.Add(
		leaflet.NewTileLayer(
			leaflet.TileLayerOptions{
				MaxZoom:     18,
				Attribution: `&copy; <a href="http://www.openstreetmap.org/copyright">OpenStreetMap</a>`,
			},
		),
	)

	p.updateModel()
	p.updateMapFromModel()
}

// Render implements the vecty.Component interface.
func (p *Home) Render() vecty.ComponentOrHTML {
	hasElement := p.e.entity != nil
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
					p.m,
				),
				vecty.If(hasElement,
					elem.Div(
						md.LayoutGridCell(
							md.LayoutGridCellOptions{Span: 4},
							&p.e,
						),
					),
				),
				md.LayoutGridCell(md.LayoutGridCellOptions{Span: 12},
					p.slider,
				),
			),
		),
	)
}
