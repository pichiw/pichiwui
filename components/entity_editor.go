package components

import (
	"github.com/gowasm/vecty"
	"github.com/gowasm/vecty/elem"
)

func NewEntityEditor(e *Entity) *EntityEditor {
	return &EntityEditor{entity: e}
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
