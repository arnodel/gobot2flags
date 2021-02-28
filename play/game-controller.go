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
	selectedControl GameControl
	icons           sprites.Icons
}

func (g *gameControlSelector) Bounds() image.Rectangle {
	return image.Rect(0, 0, 24*len(gameControls), 32)
}

func (g *gameControlSelector) Draw(c engine.Canvas) {
	for i, gc := range gameControls {
		img := g.icons.Get(gameControlIcons[i])
		var opts ebiten.DrawImageOptions
		g.icons.Anchor(&opts.GeoM)
		if g.selectedControl != gc {
			opts.GeoM.Scale(0.5, 0.5)
		}
		opts.GeoM.Translate((float64(i)+0.5)*24, 16)
		c.DrawImage(img, &opts)
	}
}

func (g *gameControlSelector) Click(x, y float64) {
	idx := int(x / 24)
	if idx >= 0 && idx < len(gameControls) {
		g.selectedControl = gameControls[idx]
	}
}
