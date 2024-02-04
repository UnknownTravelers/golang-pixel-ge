package colliders

import (
	"math"

	"github.com/faiface/pixel"
)

type line pixel.Line

func (l line) Contains(col Collider) *CollisionInfo {
	switch col.(type) {
	case rect:
		return l.Rect(col.(rect))
	case circle:
		return l.Circle(col.(circle))
	case line:
		return l.Line(col.(line))
	case vec:
		return l.Vec(col.(vec))
	default:
		return col.Contains(l)
	}
}

func (l line) Rect(r rect) *CollisionInfo {
	colinfo := r.Line(l)
	if colinfo != nil {
		colinfo.Normal = colinfo.Normal.Scaled(-1)
		return colinfo
	}
	return nil
}

func (l line) Circle(c circle) *CollisionInfo {
	// Get the point on the line closest to the center of the circle.
	closest := l.Closest(vec(c.Center))
	cirToClosest := vec(c.Center).To(closest)

	if cirToClosest.Len() >= c.Radius {
		return nil
	}
	return &CollisionInfo{
		Point:  closest,
		Normal: cirToClosest.Normalized(),
	}
}

func (l line) Line(k line) *CollisionInfo {
	lm, lb := pixel.Line(l).Formula()
	km, kb := pixel.Line(k).Formula()

	if lm == km { // same slope
		if lb != kb {
			return nil
		} // same infinite line
		a := l.A.X
		b := l.B.X
		y := k.A.X
		z := k.B.X
		if a <= y && y <= b || b <= y && y <= a { // 'start' of k inside l
			return &CollisionInfo{
				Point:  vec(k.A),
				Normal: vec(pixel.ZV),
			}
		}
		if a <= z && z <= b || b <= z && z <= a { // 'end' of k inside l
			return &CollisionInfo{
				Point:  vec(k.B),
				Normal: vec(pixel.ZV),
			}
		}
		if y <= a && a <= z || z <= a && a <= y { // 'start of l inside k
			return &CollisionInfo{
				Point:  vec(l.A),
				Normal: vec(pixel.ZV),
			}
		}
		if y <= b && b <= z || z <= b && b <= y { // 'end of l inside k
			return &CollisionInfo{
				Point:  vec(l.B),
				Normal: vec(pixel.ZV),
			}
		}
	}

	var x, y float64

	if math.IsInf(math.Abs(lm), 1) || math.IsInf(math.Abs(km), 1) {
		// One line is vertical
		intersectM := lm
		intersectB := lb
		verticalLine := k

		if math.IsInf(math.Abs(lm), 1) {
			intersectM = km
			intersectB = kb
			verticalLine = l
		}

		y = intersectM*verticalLine.A.X + intersectB
		x = verticalLine.A.X
	} else {
		// Coordinates of intersect
		x = (kb - lb) / (lm - km)
		y = lm*x + lb
	}

	// if point is on both lines, they intersect
	if pixel.Line(l).Contains(pixel.V(x, y)) && pixel.Line(k).Contains(pixel.V(x, y)) {
		return &CollisionInfo{
			Point:  V(x, y),
			Normal: vec(pixel.ZV),
		}
	}
	return nil
}

func (l line) Vec(v vec) *CollisionInfo {
	if pixel.Line(l).Contains(pixel.Vec(v)) {
		return &CollisionInfo{
			Point:  v,
			Normal: l.Slope().Normal(),
		}
	}
	return nil
}

func (l line) Slope() vec {
	return vec(l.A.To(l.B))
}

func (l line) Bounds() rect {
	return rect(pixel.Line(l).Bounds())
}

func (l line) Closest(v vec) vec {
	return vec(pixel.Line(l).Closest(pixel.Vec(v)))
}

func L(from, to pixel.Vec) line {
	return line(pixel.L(from, to))
}
