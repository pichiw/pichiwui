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
	vecty.RenderBody(&Home{
		m: m,
	})

	go func() {
		for {
			for _, c := range coordinates {
				time.Sleep(time.Second * 3)
				m.View(c, m.Zoom())
			}
		}
	}()

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
