package main

import (
	"github.com/gowasm/vecty"
	"github.com/gowasm/vecty/elem"
	"github.com/pichiw/leaflet"
)

func main() {
	c := make(chan struct{}, 0)

	vecty.SetTitle("Pichiw")
	vecty.AddStylesheet("app.css")

	m := leaflet.NewMap(
		"mapid",
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
			leaflet.NewCoordinate(49.8951, -97.1384),
			leaflet.NewCoordinate(48.8951, -97.1384),
		),
	)

	vecty.RenderBody(&Home{
		m: m,
	})

	<-c
}

// Home is our main page component.
type Home struct {
	vecty.Core
	m *leaflet.Map
}

// Render implements the vecty.Component interface.
func (p *Home) Render() vecty.ComponentOrHTML {
	return elem.Body(
		elem.Div(
			p.m,
		),
	)
}
