package selectlevel

import (
	"image"
	"image/color"

	"github.com/arnodel/gobot2flags/engine"
	"github.com/hajimehoshi/ebiten/v2"
)

type View struct {
	levels []string
	// pointer       *engine.PointerTracker
	selectedLevel int
	selectLevel   func(int)
}

var _ engine.View = (*View)(nil)

func NewView(levels []string, selectLevel func(int)) *View {
	return &View{
		levels: levels,
		// pointer:       &engine.PointerTracker{},
		selectLevel:   selectLevel,
		selectedLevel: -1,
	}
}

func (v *View) Update(vc engine.ViewContainer) error {
	pointer := vc.Pointer()
	switch pointer.Status() {
	case engine.TouchDown:
		v.selectedLevel = v.selectingLevel(pointer.CurrentPos())
	case engine.Dragging:
		i := v.selectingLevel(pointer.CurrentPos())
		if i != v.selectedLevel {
			v.selectedLevel = -1
		}
	case engine.TouchUp:
		// Play it!
		if v.selectedLevel != -1 {
			lvl := v.selectedLevel
			v.selectedLevel = -1
			v.selectLevel(lvl)
		}
	}
	return nil
}

func (v *View) Draw(screen *ebiten.Image) {
	for i, level := range v.levels {
		var col color.Color
		if i == v.selectedLevel {
			col = color.RGBA{0xff, 0, 0, 0xff}
		} else {
			col = color.White
		}
		engine.DrawText(screen, level, 10, 20+i*30, col)
	}
}

func (v *View) selectingLevel(pos image.Point) int {
	for i, level := range v.levels {
		if pos.In(engine.TextBounds(level, 10, 20+i*30)) {
			return i
		}
	}
	return -1
}
