package objects

import (
	"github.com/faiface/pixel"
	"github.com/faiface/pixel/imdraw"
	"github.com/unknownTravelers/golang-pixel-ge/colliders"
)

type goal struct {
	pos    pixel.Vec
	radius float64
	step   float64

	counter float64
	cols    [5]pixel.RGBA
}

func (g *goal) Update(dt float64) {
	g.counter += dt
	for g.counter > g.step {
		g.counter -= g.step
		for i := len(g.cols) - 2; i >= 0; i-- {
			g.cols[i+1] = g.cols[i]
		}
		g.cols[0] = RandomNiceColor()
	}
}

func (g *goal) Draw(imd *imdraw.IMDraw) {
	for i := len(g.cols) - 1; i >= 0; i-- {
		imd.Color = g.cols[i]
		imd.Push(g.pos)
		imd.Circle(float64(i+1)*g.radius/float64(len(g.cols)), 0)
	}
}

func (g *goal) Collide(col colliders.Collider) *colliders.CollisionInfo {
	return nil
}

func NewGoal(pos pixel.Vec, rad, step float64) *goal {
	return &goal{
		pos:    pos,
		radius: rad,
		step:   step,
	}
}
