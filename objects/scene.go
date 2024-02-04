package objects

import "github.com/faiface/pixel/imdraw"

type scene struct {
	objects []Object
}

func (s *scene) AddObjects(o ...Object) {
	s.objects = append(s.objects, o...)
}

func (s *scene) Update(dt float64) {
	for _, obj := range s.objects {
		obj.Update(dt)
	}
}
func (s *scene) Draw(imd *imdraw.IMDraw) {
	for _, obj := range s.objects {
		obj.Draw(imd)
	}
}

func NewScene() *scene {
	return &scene{
		objects: make([]Object, 0),
	}
}
