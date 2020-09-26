package main

import "github.com/hajimehoshi/ebiten"

type Canvas interface {
	DrawImage(img *ebiten.Image, options *ebiten.DrawImageOptions) error
	Draw(ImageToDraw) error
}

type transformCanvas struct {
	target   *ebiten.Image
	baseGeoM ebiten.GeoM
}

func (s *transformCanvas) DrawImage(img *ebiten.Image, options *ebiten.DrawImageOptions) error {
	op := *options
	op.GeoM.Concat(s.baseGeoM)
	return s.target.DrawImage(img, &op)
}

func (s *transformCanvas) Draw(toDraw ImageToDraw) error {
	return s.DrawImage(toDraw.Image, toDraw.Options)
}
