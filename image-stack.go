package main

import (
	"sort"

	"github.com/hajimehoshi/ebiten"
)

type ImageStack struct {
	toDraw []ImageToDraw
}

func (s *ImageStack) Add(toDraw ImageToDraw) {
	s.toDraw = append(s.toDraw, toDraw)
}

func (s *ImageStack) Draw(c Canvas) {
	imgs := s.toDraw
	sort.Slice(imgs, func(i, j int) bool {
		return imgs[i].Z < imgs[j].Z
	})
	for _, toDraw := range imgs {
		c.DrawImage(toDraw.Image, toDraw.Options)
	}
}

func (s *ImageStack) Empty() {
	s.toDraw = s.toDraw[:0]
}

type ImageToDraw struct {
	Image   *ebiten.Image
	Options *ebiten.DrawImageOptions
	Z       float64
}
