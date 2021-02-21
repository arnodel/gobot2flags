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

func (w *Window) Contains(pt image.Point) bool {
	return pt.In(w.bounds)
}

func (w *Window) Coords(pt image.Point) (float64, float64) {
	inv := w.tr
	inv.Invert()
	return inv.Apply(float64(pt.X), float64(pt.Y))
}

func centeredWindow(bounds image.Rectangle, drawBounds image.Rectangle, tr ebiten.GeoM) *Window {
	scaledDrawBounds := image.Rectangle{
		Min: pt(tr.Apply(fc(drawBounds.Min))),
		Max: pt(tr.Apply(fc(drawBounds.Max))),
	}
	sz := scaledDrawBounds.Size()
	diff := bounds.Size().Sub(sz).Div(2)
	wbounds := image.Rectangle{Min: bounds.Min.Add(diff), Max: bounds.Max.Sub(diff)}
	translateVec := wbounds.Min.Sub(scaledDrawBounds.Min)
	tr.Translate(fc(translateVec))
	return &Window{
		bounds: wbounds,
		tr:     tr,
	}
}

func fc(pt image.Point) (float64, float64) {
	return float64(pt.X), float64(pt.Y)
}

func pt(x, y float64) image.Point {
	return image.Pt(int(x), int(y))
}
