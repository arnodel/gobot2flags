package engine

import (
	"image"
	"math"
)

func RectFitScale(outer, inner image.Rectangle) (f float64) {
	ow, oh := fc(outer.Size())
	iw, ih := fc(inner.Size())
	return math.Min(1, math.Min(ow/iw, oh/ih))
}

// CenterRect returns the coordinates of the translation
func CenterRect(outer, inner image.Rectangle) image.Point {
	diff := outer.Size().Sub(inner.Size()).Div(2)
	return outer.Min.Add(diff).Sub(inner.Min)
}
