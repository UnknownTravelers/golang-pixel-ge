package controls

import (
	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"
)

var Controls = pixel.ZV

func Update(win *pixelgl.Window) {
	Controls = pixel.ZV
	if win.Pressed(pixelgl.KeyLeft) {
		Controls.X--
	}
	if win.Pressed(pixelgl.KeyRight) {
		Controls.X++
	}
	if win.JustPressed(pixelgl.KeyUp) {
		Controls.Y = 1
	}
}
