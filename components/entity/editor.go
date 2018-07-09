package entity

import (
	"sync/atomic"

	"github.com/gowasm/vecty"
	"github.com/gowasm/vecty/elem"
	"github.com/pichiw/pichiwui/model"
)

func NewEditor() *Editor {
	return &Editor{}
}

type Editor struct {
	vecty.Core
	entity atomic.Value
}

func (p *Editor) SetEntity(e *model.Entity) {
	p.entity.Store(e)
}

func (p *Editor) Entity() *model.Entity {
	e, ok := p.entity.Load().(*model.Entity)
	if !ok {
		return nil
	}
	return e
}

func (p *Editor) Render() vecty.ComponentOrHTML {
	e, ok := p.entity.Load().(*model.Entity)
	if !ok || e == nil {
		return elem.NoScript()
	}

	return elem.Div(
		vecty.Markup(
			vecty.Attribute("id", "entity-editor-"+e.Name),
		),
		elem.Heading2(
			vecty.Text(e.Name),
		),
	)
}
