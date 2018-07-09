package main

import (
	"github.com/gowasm/vecty"
	"github.com/gowasm/vecty/elem"
	"github.com/pichiw/leaflet"
	"github.com/pichiw/pichiwui/components/entity"
	"github.com/pichiw/pichiwui/model"
)

func main() {
	c := make(chan struct{}, 0)

	entities := []*model.Entity{
		&model.Entity{
			Name:  "Red River",
			Coord: leaflet.NewCoordinate(49.8951, -97.1384),
		},
		&model.Entity{
			Name:  "Turtle Mountain",
			Coord: leaflet.NewCoordinate(48.8469, -99.8011),
		},
		&model.Entity{
			Name:  "St. Paul des Metis",
			Coord: leaflet.NewCoordinate(53.8896, -111.4657),
		},
		&model.Entity{
			Name:  "Vancouver",
			Coord: leaflet.NewCoordinate(49.2827, -123.1207),
		},
	}

	b := &Home{
		entities: entity.NewList(entities),
	}

	vecty.SetTitle("Pichiw")

	vecty.RenderBody(b)

	<-c
}

type Home struct {
	vecty.Core

	entities *entity.List
}

// Render implements the vecty.Component interface.
func (p *Home) Render() vecty.ComponentOrHTML {
	return elem.Body(p.entities)
}
