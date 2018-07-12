package model

// Perspective stores a chain of perspectives
type Perspective struct {
	Entities []*Entity
	Children []*Perspective
}

func (p *Perspective) AllEntities() []*Entity {
	entities := make([]*Entity, len(p.Entities))
	copy(entities, p.Entities)
	for _, c := range p.Children {
		entities = append(entities, c.AllEntities()...)
	}
	return entities
}
