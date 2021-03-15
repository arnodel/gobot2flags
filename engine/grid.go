package engine

import "image"

type Grid struct {
	CellWidth, CellHeight float64
	Columns, Rows         int
}

func (g *Grid) Bounds() image.Rectangle {
	return image.Rect(0, 0, int(g.CellWidth*float64(g.Columns)), int(g.CellHeight*float64(g.Rows)))
}

func (g *Grid) CellCenter(x, y int) (float64, float64) {
	y += x / g.Columns
	x %= g.Columns
	return (float64(x) + 0.5) * g.CellWidth, (float64(y) + 0.5) * g.CellHeight
}

func (g *Grid) CellCoords(x, y float64) (int, int) {
	x /= g.CellWidth
	y /= g.CellHeight
	i, j := int(x), int(y)
	if i < 0 || j < 0 || i >= g.Columns || j >= g.Rows {
		return -1, 0
	}
	return i, j
}

func (g *Grid) CellIndex(x, y float64) int {
	i, j := g.CellCoords(x, y)
	return i + j*g.Columns
}

func (g *Grid) CellBounds(x, y int) image.Rectangle {
	y += x / g.Columns
	x %= g.Columns
	xx, yy := g.CellWidth*float64(x), g.CellHeight*float64(y)

	return image.Rectangle{Min: pt(xx, yy), Max: pt(xx+g.CellWidth, yy+g.CellHeight)}
}
