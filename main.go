package main

import (
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

	b := &Home{
		model: []*Entity{
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
		},
		m: leaflet.NewMap(
			"mapid",
			leaflet.MapOptions{
				Center:  leaflet.NewCoordinate(49.8951, -97.1384),
				Zoom:    6,
				MaxZoom: 18,
			},
			nil,
			leaflet.NewTileLayer(
				leaflet.TileLayerOptions{
					MaxZoom:     18,
					Attribution: `&copy; <a href="http://www.openstreetmap.org/copyright">OpenStreetMap</a>`,
				},
			),
		),
	}

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

	model []*Entity
}

func (p *Home) updateMapFromModel() {
	colors := []string{
		"#ff0000",
		"#aa2200",
		"#66aa00",
		"#00ff00",
	}

	for i := range p.model {
		entity := p.model[i]

		if i > 0 {
			p.m.Add(
				leaflet.NewPolyline(
					leaflet.PolylineOptions{
						PathOptions: leaflet.PathOptions{
							Color: colors[i],
						},
					},
					p.model[i-1].Coord,
					entity.Coord,
				),
			)
		}

		p.m.Add(
			leaflet.NewMarker(
				entity.Coord,
				leaflet.Events{
					"click": func(vs []js.Value) {
						if p.e.entity == nil || p.e.entity != entity {
							p.e.entity = entity
							p.m.View(entity.Coord, p.m.Zoom())
						} else {
							p.e.entity = nil
						}
						vecty.Rerender(p)
					},
				},
			),
		)
	}
}

// Render implements the vecty.Component interface.
func (p *Home) Render() vecty.ComponentOrHTML {
	p.updateMapFromModel()
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
						md.LayoutGridCell(md.LayoutGridCellOptions{Span: 4},
							&p.e,
						),
					),
				),
			),
		),
	)
}
