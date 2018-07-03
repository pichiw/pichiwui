package main

import (
	"github.com/gowasm/gopherwasm/js"

	"github.com/gowasm/vecty"
	"github.com/gowasm/vecty/elem"
	"github.com/pichiw/leaflet"
)

func main() {
	c := make(chan struct{}, 0)

	vecty.SetTitle("Pichiw")
	vecty.AddStylesheet("app.css")

	entities := []*Entity{
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

	var coordinates []*leaflet.Coordinate

	for _, e := range entities {
		coordinates = append(coordinates, e.Coord)
	}

	m := leaflet.NewMap(
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
	)

	b := &Home{
		m: m,
	}

	colors := []string{
		"#ff0000",
		"#aa2200",
		"#66aa00",
		"#00ff00",
	}

	for i := range entities {
		entity := entities[i]

		if i > 0 {
			m.Add(
				leaflet.NewPolyline(
					leaflet.PolylineOptions{
						PathOptions: leaflet.PathOptions{
							Color: colors[i],
						},
					},
					entities[i-1].Coord,
					entity.Coord,
				),
			)
		}

		m.Add(
			leaflet.NewMarker(
				entity.Coord,
				leaflet.Events{
					"click": func(vs []js.Value) {
						if b.e.entity == nil || b.e.entity != entity {
							b.e.entity = entity
							b.m.View(entity.Coord, b.m.Zoom())
						} else {
							b.e.entity = nil
						}
						vecty.Rerender(b)
					},
				},
			),
		)
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

// Home is our main page component.
type Home struct {
	vecty.Core
	m *leaflet.Map
	e EntityEditor
}

// Render implements the vecty.Component interface.
func (p *Home) Render() vecty.ComponentOrHTML {
	hasElement := p.e.entity != nil
	return elem.Body(
		elem.Div(
			vecty.Markup(
				vecty.Class("mdc-layout-grid"),
			),
			elem.Div(
				vecty.Markup(
					vecty.Class("mdc-layout-grid__inner"),
				),
				elem.Div(
					vecty.Markup(
						vecty.Class("mdc-layout-grid__cell", "mdc-layout-grid__cell--span-2"),
					),
					elem.Heading1(
						vecty.Text("Pichiw"),
					),
				),
				elem.Div(
					vecty.Markup(
						vecty.MarkupIf(hasElement,
							vecty.Class("mdc-layout-grid__cell", "mdc-layout-grid__cell--span-6"),
						),
						vecty.MarkupIf(!hasElement,
							vecty.Class("mdc-layout-grid__cell", "mdc-layout-grid__cell--span-10"),
						),
					),
					p.m,
				),
				vecty.If(
					hasElement,
					elem.Div(
						vecty.Markup(
							vecty.Class("mdc-layout-grid__cell", "mdc-layout-grid__cell--span-4"),
						),
						&p.e,
					),
				),
			),
		),
	)
}
