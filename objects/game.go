package objects

type game struct {
	scenes       []*scene
	currentScene *scene
}

var Game *game

func init() {
	Game = &game{
		scenes:       make([]*scene, 0),
		currentScene: nil,
	}
}

func (g *game) GetCurrentScene() *scene {
	return g.currentScene
}

func (g *game) AddScenes(s ...*scene) {
	g.scenes = append(g.scenes, s...)
	if g.currentScene == nil {
		g.currentScene = s[0]
	}
}
