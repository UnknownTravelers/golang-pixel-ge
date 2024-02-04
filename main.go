package main

import (
	"fmt"
	"math"
	"time"

	_ "image/png"

	"github.com/faiface/pixel"
	"github.com/faiface/pixel/imdraw"
	"github.com/faiface/pixel/pixelgl"
	"github.com/faiface/pixel/text"
	"github.com/unknownTravelers/golang-pixel-ge/controls"
	"github.com/unknownTravelers/golang-pixel-ge/loader"
	"github.com/unknownTravelers/golang-pixel-ge/objects"
	"golang.org/x/image/colornames"
)

const (
	width  = 1024
	height = 768
)

func run() {
	sheet, anims, err := loader.AnimationSheet("sheet.png", "sheet.csv", 12)
	if err != nil {
		panic(err)
	}

	cfg := pixelgl.WindowConfig{
		Title:  "Platformer",
		Bounds: pixel.R(0, 0, width, height),
		VSync:  true,
	}
	win, err := pixelgl.NewWindow(cfg)
	if err != nil {
		panic(err)
	}

	// Create level
	scene := objects.NewScene()

	// hardcoded level
	platforms := []objects.Object{
		objects.NewPlatform(pixel.R(-50, -34, 50, -32)),
		objects.NewPlatform(pixel.R(20, 0, 70, 2)),
		objects.NewPlatform(pixel.R(-100, 10, -50, 12)),
		objects.NewPlatform(pixel.R(120, -22, 140, -20)),
		objects.NewPlatform(pixel.R(120, -72, 140, -70)),
		objects.NewPlatform(pixel.R(120, -122, 140, -120)),
		objects.NewPlatform(pixel.R(-100, -152, 100, -150)),
		objects.NewPlatform(pixel.R(-150, -127, -140, -125)),
		objects.NewPlatform(pixel.R(-180, -97, -170, -95)),
		objects.NewPlatform(pixel.R(-150, -67, -140, -65)),
		objects.NewPlatform(pixel.R(-180, -37, -170, -35)),
		objects.NewPlatform(pixel.R(-150, -7, -140, -5)),
	}
	gol := objects.NewGoal(pixel.V(-75, 40), 18, 1.0/7)

	scene.AddObjects(platforms...)
	scene.AddObjects(gol)

	// Creating player & add it to level
	goph := objects.NewGopher(sheet, anims)
	scene.AddObjects(goph)

	objects.Game.AddScenes(scene)

	// Creating window
	canvas := pixelgl.NewCanvas(pixel.R(-160/2, -120/2, 160/2, 120/2))
	imd := imdraw.New(sheet)
	imd.Precision = 32

	camPos := pixel.ZV

	last := time.Now()
	for !win.Closed() {
		dt := time.Since(last).Seconds()
		last = time.Now()

		// lerp the camera position towards the gopher
		camPos = pixel.Lerp(camPos, goph.Phys.Pos, 1-math.Pow(1.0/128, dt))
		cam := pixel.IM.Moved(camPos.Scaled(-1))
		canvas.SetMatrix(cam)

		// slow motion with tab
		if win.Pressed(pixelgl.KeyTab) {
			dt /= 8
		}

		// restart the level on pressing enter
		if win.JustPressed(pixelgl.KeyEnter) {
			goph.Phys.Rect = goph.Phys.Rect.Moved(goph.Phys.Rect.Center().Scaled(-1))
			goph.Phys.Vel = pixel.ZV
		}

		// control the gopher with keys
		controls.Update(win)

		// update the physics and animation
		scene.Update(dt)

		// draw the scene to the canvas using IMDraw
		canvas.Clear(colornames.Black)
		imd.Clear()
		scene.Draw(imd)
		imd.Draw(canvas)

		// stretch the canvas to the window
		win.Clear(colornames.White)
		win.SetMatrix(pixel.IM.Scaled(pixel.ZV,
			math.Min(
				win.Bounds().W()/canvas.Bounds().W(),
				win.Bounds().H()/canvas.Bounds().H(),
			),
		).Moved(win.Bounds().Center()))
		canvas.Draw(win, pixel.IM.Moved(canvas.Bounds().Center()))
		// Debug text
		debug_info := text.New(pixel.ZV, text.Atlas7x13)
		debug_info.WriteString(fmt.Sprintf("%v", 1/dt))
		debug_info.Draw(win, pixel.IM.Scaled(pixel.ZV, 0.55).Moved(canvas.Bounds().Min))

		win.Update()
	}
}

func main() {
	pixelgl.Run(run)
}
