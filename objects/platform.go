package objects

import (
	"image/color"

	"github.com/faiface/pixel"
	"github.com/faiface/pixel/imdraw"
	"github.com/unknownTravelers/3D-jump-infinite/colliders"
)

type platform struct {
	Rect  pixel.Rect
	Color color.Color
}

func (p *platform) Draw(imd *imdraw.IMDraw) {
	imd.Color = p.Color
	imd.Push(p.Rect.Min, p.Rect.Max)
	imd.Rectangle(0)
}

func (p *platform) Update(dt float64) {}

func (p *platform) Collide(col colliders.Collider) *colliders.CollisionInfo {
	return colliders.Rect(p.Rect).Contains(col)
}

func NewPlatform(r pixel.Rect) *platform {
	return &platform{
		Rect:  r,
		Color: RandomNiceColor(),
	}
}
