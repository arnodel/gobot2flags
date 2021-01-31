package main

import (
	"image"

	"github.com/hajimehoshi/ebiten/v2"
)

type Window struct {
	bounds image.Rectangle
	tr     ebiten.GeoM
}

func (w *Window) Center() {
	bounds := w.bounds
	w.tr.Translate(float64(bounds.Min.X+bounds.Max.X)/2, float64(bounds.Min.Y+bounds.Max.Y)/2)
}

func (w *Window) Canvas(screen *ebiten.Image) Canvas {
	return &transformCanvas{
		target:   screen.SubImage(w.bounds).(*ebiten.Image),
		baseGeoM: w.tr,
	}
}

func (w *Window) Contains(x, y int) bool {
	return image.Pt(x, y).In(w.bounds)
}

func (w *Window) Coords(x, y int) (float64, float64) {
	inv := w.tr
	inv.Invert()
	return inv.Apply(float64(x), float64(y))
}
