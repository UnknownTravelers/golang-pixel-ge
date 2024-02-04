package colliders

type merged struct {
	colliders []Collider
}

func (m *merged) Contains(col Collider) *CollisionInfo {
	for _, c := range m.colliders {
		colinfo := c.Contains(col)
		if colinfo != nil {
			return colinfo // NOTE: think about getting all colinfo and merge them (somehow)
		}
	}
	return nil
}

func (c *merged) AddColliders(col ...Collider) {
	c.colliders = append(c.colliders, col...)
}

func M(cols ...Collider) *merged {
	m := &merged{
		colliders: make([]Collider, 0),
	}
	m.AddColliders(cols...)
	return m
}
