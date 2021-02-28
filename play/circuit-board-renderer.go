package play

import (
	"image"

	"github.com/arnodel/gobot2flags/engine"
	"github.com/arnodel/gobot2flags/model"
	"github.com/hajimehoshi/ebiten/v2"
)

type ChipRenderer struct {
	*engine.Sprite
}

func (r ChipRenderer) GetChipImage(c model.ChipType) *ebiten.Image {
	return r.GetImage(chipType2imageIdx[c], 0)
}

func (r ChipRenderer) GetArrowImage(a model.ArrowType) *ebiten.Image {
	return r.GetImage(arrowType2ImageIdx[a], 0)
}

type CircuitBoardRenderer struct {
	chips         ChipRenderer
	width, height float64
}

func NewCircuitBoardRenderer(chips ChipRenderer) CircuitBoardRenderer {
	return CircuitBoardRenderer{
		chips:  chips,
		width:  32,
		height: 32,
	}
}

func (r *CircuitBoardRenderer) DrawCircuitBoard(c engine.Canvas, b *model.CircuitBoard) {
	w, h := b.Size()

	// Draw the background
	for y := 0; y < h; y++ {
		for x := 0; x < w; x++ {
			c.Draw(r.Background(x, y, false))
		}
	}

	// Draw the arrows
	for y := 0; y < h; y++ {
		for x := 0; x < w; x++ {
			chip := b.ChipAt(x, y)
			if o, ok := chip.ArrowYes(); ok {
				if chip.IsTest() {
					c.Draw(r.ArrowYes(x, y, o, chip.IsArrowActive(o)))
				} else {
					c.Draw(r.Arrow(x, y, o, chip.IsArrowActive(o)))
				}
			}
			if o, ok := chip.ArrowNo(); ok && chip.IsTest() {
				c.Draw(r.ArrowNo(x, y, o, chip.IsArrowActive(o)))
			}
		}
	}

	// Draw the chips
	for y := 0; y < h; y++ {
		for x := 0; x < w; x++ {
			chip := b.ChipAt(x, y)
			if chip.Type() != model.NoChip {
				c.Draw(r.Chip(chip.Type(), x, y, chip.IsActive()))
			}
		}
	}
}

func (r *CircuitBoardRenderer) CircuitBoardBounds(b *model.CircuitBoard) image.Rectangle {
	w, h := b.Size()
	return image.Rect(0, 0, w*32, h*32)
}

func (b CircuitBoardRenderer) GetSlotCoords(x, y float64) (int, int) {
	return int(x / b.width), int(y / b.height)
}

func (b CircuitBoardRenderer) Chip(c model.ChipType, x, y int, active bool) engine.ImageToDraw {
	return b.imageByIndex(chipType2imageIdx[c], x, y, active)
}

func (b CircuitBoardRenderer) Background(x, y int, active bool) engine.ImageToDraw {
	return b.imageByIndex(backgroundIdx, x, y, active)
}

func (b CircuitBoardRenderer) Arrow(x, y int, o model.Orientation, active bool) engine.ImageToDraw {
	return b.arrow(arrowNorthIdx, x, y, o, active)
}

func (b CircuitBoardRenderer) ArrowYes(x, y int, o model.Orientation, active bool) engine.ImageToDraw {
	return b.arrow(arrowYesNorthIdx, x, y, o, active)
}

func (b CircuitBoardRenderer) ArrowNo(x, y int, o model.Orientation, active bool) engine.ImageToDraw {
	return b.arrow(arrowNoNorthIdx, x, y, o, active)
}

func (b CircuitBoardRenderer) imageByIndex(i int, x, y int, active bool) engine.ImageToDraw {
	opts := ebiten.DrawImageOptions{}
	b.chips.Anchor(&opts.GeoM)
	opts.GeoM.Translate((float64(x)+0.5)*b.width, (float64(y)+0.5)*b.height)
	return engine.ImageToDraw{
		Image:   b.chips.GetImage(i, activeFrame(active)),
		Options: &opts,
		Z:       chipZ,
	}
}

func (b CircuitBoardRenderer) arrow(baseIdx int, x, y int, o model.Orientation, active bool) engine.ImageToDraw {
	// log.Printf("Arrow %d, x=%d, y=%d, o=%d", baseIdx, x, y, o)
	i := baseIdx + int(o)
	v := o.VelocityForward()
	opts := ebiten.DrawImageOptions{}
	b.chips.Anchor(&opts.GeoM)
	opts.GeoM.Translate(
		(float64(x)+0.5*(1+float64(v.Dx)))*b.width,
		(float64(y)+0.5*(1+float64(v.Dy)))*b.height,
	)
	return engine.ImageToDraw{
		Image:   b.chips.GetImage(i, activeFrame(active)),
		Options: &opts,
		Z:       arrowZ,
	}
}

func activeFrame(active bool) int {
	if active {
		return 1
	}
	return 0
}

var chipType2imageIdx = map[model.ChipType]int{
	model.StartChip:         startIdx,
	model.ForwardChip:       forwardIdx,
	model.TurnLeftChip:      turnLeftIdx,
	model.TurnRightChip:     turnRightIdx,
	model.PaintRedChip:      paintRedIdx,
	model.PaintYellowChip:   paintYellowIdx,
	model.PaintBlueChip:     paintBlueIdx,
	model.IsWallAheadChip:   isWallAheadIdx,
	model.IsFloorRedChip:    isFloorRedIdx,
	model.IsFloorYellowChip: isFloorYellowIdx,
	model.IsFloorBlueChip:   isFloorBueIdx,
}

var arrowType2ImageIdx = map[model.ArrowType]int{
	model.ArrowNo:  arrowNoNorthIdx,
	model.ArrowYes: arrowYesNorthIdx,
}

const (
	backgroundIdx = iota
	startIdx
	isFloorRedIdx
	isFloorYellowIdx
	isFloorBueIdx
	isWallAheadIdx
	forwardIdx
	turnRightIdx
	turnLeftIdx
	paintRedIdx
	paintYellowIdx
	paintBlueIdx

	arrowNorthIdx
	arrowEastIdx
	arrowSouthIdx
	arrowWestIdx

	arrowYesNorthIdx
	arrowYesEastIdx
	arrowYesSouthIdx
	arrowYesWestIdx

	arrowNoNorthIdx
	arrowNoEastIdx
	arrowNoSouthIdx
	arrowNoWestIdx
)

const (
	bgZ float64 = iota
	arrowZ
	chipZ
)
