package main

import (
	"time"

	"github.com/gowasm/vecty"
	"github.com/gowasm/vecty/elem"
	"github.com/pichiw/leaflet"
	"github.com/pichiw/pichiwui/components/perspective"
	"github.com/pichiw/pichiwui/model"
)

func main() {
	c := make(chan struct{}, 0)

	model := &model.Perspective{
		Entities: []*model.Entity{
			&model.Entity{
				Name:  "Red River",
				Coord: leaflet.NewCoordinate(49.8951, -97.1384),
				Time:  time.Date(1800, 6, 6, 0, 0, 0, 0, time.UTC),
			},
			&model.Entity{
				Name:  "Turtle Mountain",
				Coord: leaflet.NewCoordinate(48.8469, -99.8011),
				Time:  time.Date(1850, 6, 6, 0, 0, 0, 0, time.UTC),
			},
		},
		Children: []*model.Perspective{
			&model.Perspective{
				Entities: []*model.Entity{
					&model.Entity{
						Name:  "St. Paul des Metis",
						Coord: leaflet.NewCoordinate(53.8896, -111.4657),
						Time:  time.Date(1895, 6, 6, 0, 0, 0, 0, time.UTC),
					},
					&model.Entity{
						Name:  "Vancouver",
						Coord: leaflet.NewCoordinate(49.2827, -123.1207),
						Time:  time.Date(1980, 6, 6, 0, 0, 0, 0, time.UTC),
					},
				},
			},
			&model.Perspective{
				Entities: []*model.Entity{
					&model.Entity{
						Name:  "Regina",
						Coord: leaflet.NewCoordinate(50.4452, -104.6189),
						Time:  time.Date(1950, 6, 6, 0, 0, 0, 0, time.UTC),
					},
					&model.Entity{
						Name:  "Vancouver",
						Coord: leaflet.NewCoordinate(49.2827, -123.1207),
						Time:  time.Date(1979, 6, 6, 0, 0, 0, 0, time.UTC),
					},
				},
			},
		},
	}

	b := &Home{
		perspectiveList: perspective.NewList(model),
	}

	vecty.SetTitle("Pichiw")

	vecty.RenderBody(b)

	<-c
}

type Home struct {
	vecty.Core

	perspectiveList *perspective.List
}

// Render implements the vecty.Component interface.
func (p *Home) Render() vecty.ComponentOrHTML {
	return elem.Body(p.perspectiveList)
}
