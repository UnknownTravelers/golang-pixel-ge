package colliders

import "github.com/faiface/pixel"

type circle pixel.Circle

func (c circle) Contains(col Collider) *CollisionInfo {
	switch col.(type) {
	case rect:
		return c.Rect(col.(rect))
	case circle:
		return c.Circle(col.(circle))
	case line:
		return c.Line(col.(line))
	case vec:
		return c.Vec(col.(vec))
	default:
		return col.Contains(c)
	}
}

func (c circle) Rect(r rect) *CollisionInfo {
	// Checks if the c.Center is not in the diagonal quadrants of the rectangle
	if (r.Min.X <= c.Center.X && c.Center.X <= r.Max.X) || (r.Min.Y <= c.Center.Y && c.Center.Y <= r.Max.Y) {
		// 'grow' the Rect by c.Radius in each orthagonal
		grown := rect{Min: r.Min.Sub(pixel.V(c.Radius, c.Radius)), Max: r.Max.Add(pixel.V(c.Radius, c.Radius))}

		if pixel.Rect(grown).Contains(c.Center) {
			contact := r.ClosestPerimeter(vec(c.Center))
			return &CollisionInfo{
				Point:  vec(c.Center).To(contact).Normalized().Scaled(c.Radius),
				Normal: contact.To(vec(c.Center)).Normalized(),
			}
		}
		return nil
	} else {
		// The center is in the diagonal quadrants
		for _, v := range r.Vertices() {
			if colinfo := c.Vec(v); colinfo != nil {
				return &CollisionInfo{
					Point:  vec(c.Center).To(v).Normalized().Scaled(c.Radius),
					Normal: v.To(vec(c.Center)).Normalized(),
				}
			}
		}
		return nil
	}
}

func (c1 circle) Circle(c2 circle) *CollisionInfo {
	d := c1.Center.To(c2.Center)
	if d.Len() > c1.Radius+c2.Radius {
		return nil
	}
	return &CollisionInfo{
		Point:  vec(d.Add(c1.Center)),
		Normal: vec(d).Normalized(),
	}
}

func (c circle) Line(l line) *CollisionInfo {
	// Get the point on the line closest to the center of the circle.
	closest := l.Closest(vec(c.Center))
	cirToClosest := vec(c.Center).To(closest)

	if cirToClosest.Len() >= c.Radius {
		return nil
	}
	return &CollisionInfo{
		Point:  cirToClosest.Normalized().Scaled(c.Radius).Add(vec(c.Center)),
		Normal: cirToClosest.Normalized(),
	}
}

func (c circle) Vec(v vec) *CollisionInfo {
	if pixel.Circle(c).Contains(pixel.Vec(v)) {
		return &CollisionInfo{
			Point:  v,
			Normal: vec(c.Center).To(v).Normalized(),
		}
	}
	return nil
}

func C(cen pixel.Vec, rad float64) circle {
	return circle(pixel.C(cen, rad))
}
