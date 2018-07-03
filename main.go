package main

import (
	"time"

	"github.com/gowasm/vecty"
	"github.com/gowasm/vecty/elem"
	"github.com/pichiw/leaflet"
)

func main() {
	c := make(chan struct{}, 0)

	vecty.SetTitle("Pichiw")
	vecty.AddStylesheet("app.css")

	coordinates := leaflet.NewCoordinates(
		49.8951, -97.1384,
		48.8469, -99.8011,
		53.8896, -111.4657,
		49.2827, -123.1207,
	)

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
		leaflet.NewPolyline(
			leaflet.PolylineOptions{
				PathOptions: leaflet.PathOptions{
					Color: "#ff0000",
				},
			},
			coordinates...,
		),
	)

	for _, c := range coordinates {
		m.Add(leaflet.NewMarker(c))
	}

	e := &Element{}
	b := &Home{
		m: m,
		e: e,
	}
	vecty.RenderBody(b)

	go func() {
		time.Sleep(time.Second * 5)
		e.Name = "Hello"
		vecty.Rerender(b)
		time.Sleep(time.Second * 5)
		e.Name = ""
		vecty.Rerender(b)
	}()

	<-c
}

type Element struct {
	vecty.Core
	Name string
}

func (p *Element) Render() vecty.ComponentOrHTML {
	return elem.Div(
		elem.Heading2(
			vecty.Text(p.Name),
		),
	)
}

// Home is our main page component.
type Home struct {
	vecty.Core
	m *leaflet.Map
	e *Element
}

// Render implements the vecty.Component interface.
func (p *Home) Render() vecty.ComponentOrHTML {
	hasElement := len(p.e.Name) > 0
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
						p.e,
					),
				),
			),
		),
	)
}
