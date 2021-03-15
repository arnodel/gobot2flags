package play

import (
	"image"

	"github.com/arnodel/gobot2flags/engine"
	"github.com/arnodel/gobot2flags/sprites"
	"github.com/hajimehoshi/ebiten/v2"
)

type GameControl int

const (
	NoControl GameControl = iota
	Play
	FastForward
	Rewind
	Pause
	Step
	Exit
)

var gameControls = []GameControl{Rewind, Play, Step, Pause, FastForward}
var gameControlIcons = []sprites.IconType{sprites.RewindIcon, sprites.PlayIcon, sprites.StepIcon, sprites.PauseIcon, sprites.FastForwardIcon}

type gameControlSelector struct {
	selectedControl  GameControl
	selectingControl GameControl
	grid             engine.Grid
	indexSelector    engine.Selector
	icons            sprites.Icons
}

func (g *gameControlSelector) Bounds() image.Rectangle {
	return g.grid.Bounds()
}

func (g *gameControlSelector) Draw(c engine.Canvas) {
	for i, gc := range gameControls {
		img := g.icons.Get(gameControlIcons[i])
		var opts ebiten.DrawImageOptions
		g.icons.Anchor(&opts.GeoM)
		if g.selectedControl != gc && g.selectingControl != gc {
			opts.GeoM.Scale(0.5, 0.5)
		}
		opts.GeoM.Translate(g.grid.CellCenter(i, 0))
		c.DrawImage(img, &opts)
	}
}

func (g *gameControlSelector) Update(p engine.PointerStatus) {
	g.grid = engine.Grid{
		CellWidth:  24,
		CellHeight: 32,
		Rows:       1,
		Columns:    len(gameControls),
	}

	idx := g.grid.CellIndex(p.CurrentCoords())
	var ctrl GameControl
	if idx >= 0 && idx < len(gameControls) {
		ctrl = gameControls[idx]
	}
	switch g.indexSelector.Update(idx, p.Status()) {
	case engine.NotSelecting:
		g.selectingControl = NoControl
	case engine.Selecting:
		g.selectingControl = ctrl
	case engine.Select:
		g.selectedControl = ctrl
		g.selectingControl = NoControl
	}
}
