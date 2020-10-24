package main

import "github.com/hajimehoshi/ebiten/v2"

type Canvas interface {
	DrawImage(img *ebiten.Image, options *ebiten.DrawImageOptions)
	Draw(ImageToDraw)
}

type transformCanvas struct {
	target   *ebiten.Image
	baseGeoM ebiten.GeoM
}

var _ Canvas = &transformCanvas{}

func (s *transformCanvas) DrawImage(img *ebiten.Image, options *ebiten.DrawImageOptions) {
	op := *options
	op.GeoM.Concat(s.baseGeoM)
	s.target.DrawImage(img, &op)
}

func (s *transformCanvas) Draw(toDraw ImageToDraw) {
	s.DrawImage(toDraw.Image, toDraw.Options)
}
