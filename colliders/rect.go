package colliders

import (
	"sort"

	"github.com/faiface/pixel"
)

type rect pixel.Rect

type ByDist []vec

func (a ByDist) Len() int           { return len(a) }
func (a ByDist) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a ByDist) Less(i, j int) bool { return a[i].Len() < a[j].Len() }

func (r rect) Contains(col Collider) *CollisionInfo {
	switch col.(type) {
	case rect:
		return r.Rect(col.(rect))
	case circle:
		return r.Circle(col.(circle))
	case line:
		return r.Line(col.(line))
	case vec:
		return r.Vec(col.(vec))
	default:
		return col.Contains(r)
	}
}

func (r1 rect) Rect(r2 rect) *CollisionInfo {
	rect := pixel.Rect(r1).Intersect(pixel.Rect(r2))
	if rect == pixel.ZR {
		return nil
	}
	return &CollisionInfo{
		Point:  vec(rect.Min),
		Normal: r1.Center().To(r2.Center()).Normalized(),
	}
}

func (r rect) Circle(c circle) *CollisionInfo {
	// Checks if the c.Center is not in the diagonal quadrants of the rectangle
	if (r.Min.X <= c.Center.X && c.Center.X <= r.Max.X) || (r.Min.Y <= c.Center.Y && c.Center.Y <= r.Max.Y) {
		// 'grow' the Rect by c.Radius in each orthagonal
		grown := rect{Min: r.Min.Sub(pixel.V(c.Radius, c.Radius)), Max: r.Max.Add(pixel.V(c.Radius, c.Radius))}

		if pixel.Rect(grown).Contains(c.Center) {
			contact := r.ClosestPerimeter(vec(c.Center))
			return &CollisionInfo{
				Point:  contact,
				Normal: contact.To(vec(c.Center)).Normalized(),
			}
		}
		return nil
	} else {
		// The center is in the diagonal quadrants
		for _, v := range r.Vertices() {
			if colinfo := c.Vec(v); colinfo != nil {
				return &CollisionInfo{
					Point:  v,
					Normal: v.To(vec(c.Center)).Normalized(),
				}
			}
		}
		return nil
	}
}

func (r rect) Line(l line) *CollisionInfo {
	// Check if either end of the line segment are within the rectangle
	if pixel.Rect(r).Contains(l.A) || pixel.Rect(r).Contains(l.B) {
		return r.Rect(l.Bounds()) // NOTE: might be sufficient
	}

	// Check if any of the rectangles' edges intersect with this line.
	for _, edge := range r.Edges() {
		if colinfo := l.Line(edge); colinfo != nil {
			// TODO: Make sure colinfo.Normal is the edge Normal and it is directed outside of the rect
			colinfo.Normal = edge.Slope().Normal()
			return colinfo
		}
	}
	return nil
}

func (r rect) Vec(v vec) *CollisionInfo {
	if pixel.Rect(r).Contains(pixel.Vec(v)) {
		return &CollisionInfo{
			Point:  v,
			Normal: r.Center().To(v).Normalized(),
		}
	}
	return nil
}

func (r rect) Normalized() rect {
	if r.Min.X > r.Max.X {
		r.Min.X, r.Max.X = r.Max.X, r.Min.X
	}
	if r.Min.Y > r.Max.Y {
		r.Min.Y, r.Max.Y = r.Max.Y, r.Min.Y
	}
	return r
}

func (r rect) ClosestPerimeter(v vec) vec {
	var closests [4]vec
	for i, l := range r.Edges() {
		closests[i] = l.Closest(v)
	}
	sort.Sort(ByDist(closests[:]))
	return closests[0]
}

func (r rect) Vertices() [4]vec {
	return [4]vec{
		vec(r.Min),
		V(r.Min.X, r.Max.Y),
		vec(r.Max),
		V(r.Max.X, r.Min.Y),
	}
}

func (r rect) Edges() [4]line {
	corners := pixel.Rect(r).Vertices()

	return [4]line{
		{A: corners[0], B: corners[1]},
		{A: corners[1], B: corners[2]},
		{A: corners[2], B: corners[3]},
		{A: corners[3], B: corners[0]},
	}
}

func (r rect) Center() vec {
	return vec(pixel.Rect(r).Center())
}

func (r rect) Grow(left, bottom, right, top float64) rect {
	r.Min.X += left
	r.Min.Y += bottom
	r.Max.X += right
	r.Max.Y += top
	return r
}

func R(minX, minY, maxX, maxY float64) rect {
	return rect(pixel.R(minX, minY, maxX, maxY)).Normalized()
}

func Rect(r pixel.Rect) rect {
	return rect(r)
}
