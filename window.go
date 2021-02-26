package main

import (
	"image"
	"image/color"
	"math"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

func drawRect(dst *ebiten.Image, r image.Rectangle, c color.Color) {
	x1, y1 := fc(r.Min)
	x2, y2 := fc(r.Max)
	ebitenutil.DrawLine(dst, x1, y1, x1, y2, c)
	ebitenutil.DrawLine(dst, x1, y1, x2, y1, c)
	ebitenutil.DrawLine(dst, x2, y1, x2, y2, c)
	ebitenutil.DrawLine(dst, x2, y2, x1, y2, c)
}

type Window struct {
	outerBounds image.Rectangle
	drawBounds  image.Rectangle
	tr          ebiten.GeoM
}

func (w *Window) Center() {
	bounds := w.drawBounds
	w.tr.Translate(float64(bounds.Min.X+bounds.Max.X)/2, float64(bounds.Min.Y+bounds.Max.Y)/2)
}

func (w *Window) Canvas(screen *ebiten.Image) Canvas {
	if debug {
		drawRect(screen, w.drawBounds, color.RGBA{0xff, 0, 0, 0xff})
		drawRect(screen, w.outerBounds, color.RGBA{0, 0xff, 0, 0xff})
	}
	return &transformCanvas{
		target:   screen.SubImage(w.drawBounds).(*ebiten.Image),
		baseGeoM: w.tr,
	}
}

func (w *Window) Contains(pt image.Point) bool {
	return pt.In(w.drawBounds)
}

func (w *Window) Coords(pt image.Point) (float64, float64) {
	inv := w.tr
	inv.Invert()
	return inv.Apply(fc(pt))
}

func centeredWindow(bounds image.Rectangle, drawBounds image.Rectangle, tr ebiten.GeoM) *Window {
	minX, minY := tr.Apply(fc(drawBounds.Min))
	maxX, maxY := tr.Apply(fc(drawBounds.Max))
	bsz := bounds.Size()
	w, h := fc(bsz)
	dw, dh := maxX-minX, maxY-minY
	if w < dw || h < dh {
		scale := math.Min(w/dw, h/dh)
		tr.Scale(scale, scale)
		minX, minY = scale*minX, scale*minY
		maxX, maxY = scale*maxX, scale*maxY
	}
	scaledDrawBounds := image.Rectangle{
		Min: pt(minX, minY),
		Max: pt(maxX, maxY),
	}
	diff := bsz.Sub(scaledDrawBounds.Size()).Div(2)
	wbounds := image.Rectangle{Min: bounds.Min.Add(diff), Max: bounds.Max.Sub(diff)}
	translateVec := wbounds.Min.Sub(scaledDrawBounds.Min)
	tr.Translate(fc(translateVec))
	return &Window{
		outerBounds: bounds,
		drawBounds:  wbounds,
		tr:          tr,
	}
}

func fc(pt image.Point) (float64, float64) {
	return float64(pt.X), float64(pt.Y)
}

func pt(x, y float64) image.Point {
	return image.Pt(int(x), int(y))
}
