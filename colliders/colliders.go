// colliders wrap the pixel shapes (circle, rect, line and vec) to check if two shapes are colliding
package colliders

import (
	"errors"
)

var errUnownShapeType = errors.New("unknown collider shape")

type CollisionInfo struct {
	Point  vec
	Normal vec
}

type Collider interface {
	Contains(Collider) *CollisionInfo
}
