package selectlevel

import (
	"image"
	"image/color"

	"github.com/arnodel/gobot2flags/engine"
	"github.com/hajimehoshi/ebiten/v2"
)

type View struct {
	levels        []string
	grid          engine.Grid
	selector      engine.Selector
	selectedLevel int
	selectLevel   func(int)
}

var _ engine.View = (*View)(nil)

func NewView(levels []string, selectLevel func(int)) *View {
	return &View{
		levels:        levels,
		selectLevel:   selectLevel,
		selectedLevel: -1,
	}
}

func (v *View) Update(vc engine.ViewContainer) error {
	w, _ := vc.OutsideSize()
	v.grid = engine.Grid{
		CellWidth:  float64(w),
		CellHeight: 30,
		Columns:    1,
		Rows:       len(v.levels),
	}
	pointer := vc.Pointer()
	if v.selector.Update(v.selectingLevel(pointer.CurrentPos()), pointer.Status()) == engine.Select {
		v.selectLevel(v.selector.SelectIndex)
	}
	return nil
}

func (v *View) Draw(screen *ebiten.Image) {
	for i, level := range v.levels {
		var col color.Color
		if v.selector.IsSelecting(i) {
			col = color.RGBA{0xff, 0, 0, 0xff}
		} else {
			col = color.White
		}
		outerBox := v.grid.CellBounds(i, 0)
		textBox := engine.TextBounds(level, 0, 0)
		tr := engine.CenterRect(outerBox, textBox)
		engine.DrawText(screen, level, tr.X, tr.Y, col)
	}
}

func (v *View) selectingLevel(pos image.Point) int {
	return v.grid.CellIndex(float64(pos.X), float64(pos.Y))

}
