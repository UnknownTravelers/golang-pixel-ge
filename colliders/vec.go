package colliders

import "github.com/faiface/pixel"

type vec pixel.Vec

func (v vec) Contains(col Collider) *CollisionInfo {
	switch col.(type) {
	case circle:
		return v.Circle(col.(circle))
	case rect:
		return v.Rect(col.(rect))
	case line:
		return v.Line(col.(line))
	case vec:
		return v.Vec(col.(vec))
	default:
		return col.Contains(v)
	}
}

func (v vec) Rect(r rect) *CollisionInfo {
	if pixel.Rect(r).Contains(pixel.Vec(v)) {
		return &CollisionInfo{
			Point:  v,
			Normal: r.Center().To(v).Normalized(),
		}
	}
	return nil
}

func (v vec) Circle(c circle) *CollisionInfo {
	if pixel.Circle(c).Contains(pixel.Vec(v)) {
		return &CollisionInfo{
			Point:  v,
			Normal: vec(c.Center).To(v).Normalized(),
		}
	}
	return nil
}

func (v vec) Line(l line) *CollisionInfo {
	if pixel.Line(l).Contains(pixel.Vec(v)) {
		return &CollisionInfo{
			Point:  v,
			Normal: l.Slope().Normal(),
		}
	}
	return nil
}

func (v1 vec) Vec(v2 vec) *CollisionInfo {
	if v1.X == v2.X && v1.Y == v2.Y {
		return &CollisionInfo{
			Point:  v1,
			Normal: vec(pixel.ZV),
		}
	}
	return nil
}

func (v vec) Normalized() vec {
	l := v.Len()
	return v.Scaled(l)
}

func (v vec) Len() float64 {
	return pixel.Vec(v).Len()
}

func (v vec) Scaled(x float64) vec {
	v.X /= x
	v.Y /= x
	return v
}

func (v vec) Add(v2 vec) vec {
	v.X += v2.X
	v.Y += v2.Y
	return v
}

func (src vec) To(dst vec) vec {
	v := V(dst.X-src.X, dst.Y-src.Y)
	return v
}

func (v vec) Normal() vec {
	return vec(pixel.Vec(v).Normal())
}

func V(x, y float64) vec {
	return vec(pixel.V(x, y))
}

func Vec(v pixel.Vec) vec {
	return vec(v)
}
