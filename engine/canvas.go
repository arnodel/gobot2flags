package engine

import (
	"image"
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
)

type Canvas interface {
	DrawImage(img *ebiten.Image, options *ebiten.DrawImageOptions)
	Draw(ImageToDraw)
	Fill(color.Color)
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

func (s *transformCanvas) Clip(r image.Rectangle) *transformCanvas {
	return &transformCanvas{
		target:   s.target.SubImage(r).(*ebiten.Image),
		baseGeoM: s.baseGeoM,
	}
}

func (s *transformCanvas) Center() *transformCanvas {
	tr := s.baseGeoM
	bounds := s.target.Bounds()
	tr.Translate(float64(bounds.Min.X+bounds.Max.X)/2, float64(bounds.Min.Y+bounds.Max.Y)/2)
	return &transformCanvas{
		target:   s.target,
		baseGeoM: tr,
	}
}

func (s *transformCanvas) CursorPosition() (float64, float64) {
	x, y := ebiten.CursorPosition()
	inv := s.baseGeoM
	inv.Invert()
	return inv.Apply(float64(x), float64(y))
}

func (s *transformCanvas) Fill(c color.Color) {
	s.target.Fill(c)
}
