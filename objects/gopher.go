package objects

import (
	"math"

	"github.com/faiface/pixel"
	"github.com/faiface/pixel/imdraw"
	"github.com/unknownTravelers/golang-pixel-ge/colliders"
	"github.com/unknownTravelers/golang-pixel-ge/controls"
)

type animState int

const (
	idle animState = iota
	running
	jumping
)

type gopherAnim struct {
	sheet pixel.Picture
	anims map[string][]pixel.Rect
	rate  float64

	state   animState
	counter float64
	dir     float64

	frame pixel.Rect

	sprite *pixel.Sprite
	Phys   *gopherPhys
}

type gopherPhys struct {
	gravity   float64
	runSpeed  float64
	jumpSpeed float64

	Pos         pixel.Vec
	Rect        pixel.Rect
	GroundCheck pixel.Line
	Vel         pixel.Vec
	ground      bool
}

func (gp *gopherPhys) update(dt float64) {
	// apply controls
	switch {
	case controls.Controls.X < 0:
		gp.Vel.X = -gp.runSpeed
	case controls.Controls.X > 0:
		gp.Vel.X = +gp.runSpeed
	default:
		gp.Vel.X = 0
	}

	// apply gravity and velocity
	gp.Vel.Y += gp.gravity * dt
	gp.Pos = gp.Pos.Add(gp.Vel.Scaled(dt))

	// check collisions against each platform
	gp.ground = false
	if gp.Vel.Y <= 0 { // if falling, check ground collision
		// check if i will collide with ground between now and next frame
		r := colliders.R(gp.GroundCheck.A.X, gp.GroundCheck.A.Y, gp.GroundCheck.B.X, gp.GroundCheck.B.Y).Grow(0, -gp.Vel.Y*dt, 0, 0).Moved(colliders.Vec(gp.Pos))
		for _, o := range Game.currentScene.objects {
			// rect from line and line speed estimate
			if colinfo := o.Collide(r); colinfo != nil {
				gp.Vel.Y = 0
				gp.Pos = gp.Pos.Add(pixel.V(0, colinfo.Point.Y-r.Min.Y))
				gp.ground = true
			}
			// p, ok := o.(*platform) // TODO: needs to be generalized to object (Collide)
			// if !ok {
			// 	continue
			// }
			// if gp.Rect.Max.X <= p.Rect.Min.X || gp.Rect.Min.X >= p.Rect.Max.X {
			// 	continue
			// }
			// if gp.Rect.Min.Y > p.Rect.Max.Y || gp.Rect.Min.Y < p.Rect.Max.Y+gp.Vel.Y*dt {
			// 	continue
			// }
			// gp.Vel.Y = 0
			// gp.Rect = gp.Rect.Moved(pixel.V(0, p.Rect.Max.Y-gp.Rect.Min.Y))
			// gp.ground = true
		}
	}

	// jump if on the ground and the player wants to jump
	if gp.ground && controls.Controls.Y > 0 {
		gp.Vel.Y = gp.jumpSpeed
	}
}

func (ga *gopherAnim) Update(dt float64) {
	ga.counter += dt
	ga.Phys.update(dt)

	// determine the new animation state
	var newState animState
	switch {
	case !ga.Phys.ground:
		newState = jumping
	case ga.Phys.Vel.Len() == 0:
		newState = idle
	case ga.Phys.Vel.Len() > 0:
		newState = running
	}

	// reset the time counter if the state changed
	if ga.state != newState {
		ga.state = newState
		ga.counter = 0
	}

	// determine the correct animation frame
	switch ga.state {
	case idle:
		ga.frame = ga.anims["Front"][0]
	case running:
		i := int(math.Floor(ga.counter / ga.rate))
		ga.frame = ga.anims["Run"][i%len(ga.anims["Run"])]
	case jumping:
		speed := ga.Phys.Vel.Y
		i := int((-speed/ga.Phys.jumpSpeed + 1) / 2 * float64(len(ga.anims["Jump"])))
		if i < 0 {
			i = 0
		}
		if i >= len(ga.anims["Jump"]) {
			i = len(ga.anims["Jump"]) - 1
		}
		ga.frame = ga.anims["Jump"][i]
	}

	// set the facing direction of the gopher
	if ga.Phys.Vel.X != 0 {
		if ga.Phys.Vel.X > 0 {
			ga.dir = +1
		} else {
			ga.dir = -1
		}
	}
}

func (ga *gopherAnim) Draw(imd *imdraw.IMDraw) {
	if ga.sprite == nil {
		ga.sprite = pixel.NewSprite(nil, pixel.Rect{})
	}
	// draw the correct frame with the correct position and direction
	ga.sprite.Set(ga.sheet, ga.frame)
	ga.sprite.Draw(imd, pixel.IM.
		ScaledXY(pixel.ZV, pixel.V(
			ga.Phys.Rect.W()/ga.sprite.Frame().W(),
			ga.Phys.Rect.H()/ga.sprite.Frame().H(),
		)).
		ScaledXY(pixel.ZV, pixel.V(-ga.dir, 1)).
		Moved(ga.Phys.Pos),
	)
}

func (ga *gopherAnim) Collide(col colliders.Collider) *colliders.CollisionInfo {
	return nil
}

func NewGopher(sheet pixel.Picture, anims map[string][]pixel.Rect) *gopherAnim {
	phys := &gopherPhys{
		gravity:     -512,
		runSpeed:    64,
		jumpSpeed:   192,
		Rect:        pixel.R(-6, -7, 6, 7),
		GroundCheck: pixel.L(pixel.V(-6, -7), pixel.V(6, -7)),
	}

	anim := &gopherAnim{
		sheet: sheet,
		anims: anims,
		rate:  1.0 / 10,
		dir:   +1,
		Phys:  phys,
	}
	return anim
}
