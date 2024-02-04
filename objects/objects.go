package objects

import (
	"github.com/faiface/pixel/imdraw"
	"github.com/unknownTravelers/golang-pixel-ge/colliders"
)

type Object interface {
	Draw(*imdraw.IMDraw)
	Update(float64)
	Collide(colliders.Collider) *colliders.CollisionInfo
}
