package main

import (
	"github.com/hajimehoshi/ebiten"
)

type CircuitBoardRenderer struct {
	chips         *Sprite
	width, height float64
}

func NewBoardChipImages(img *ebiten.Image) CircuitBoardRenderer {
	return CircuitBoardRenderer{
		chips:  NewSprite(img, 32, 32, 16, 16),
		width:  32,
		height: 32,
	}
}

func (b CircuitBoardRenderer) Chip(c ChipType, x, y int) (toDraw ImageToDraw, ok bool) {
	var i int
	i, ok = chipType2imageIdx[c]
	if ok {
		toDraw = b.imageByIndex(i, x, y)
	}
	return
}

func (b CircuitBoardRenderer) Background(x, y int) ImageToDraw {
	return b.imageByIndex(backgroundIdx, x, y)
}

func (b CircuitBoardRenderer) Arrow(x, y int, o Orientation) ImageToDraw {
	return b.arrow(arrowNorthIdx, x, y, o)
}

func (b CircuitBoardRenderer) ArrowYes(x, y int, o Orientation) ImageToDraw {
	return b.arrow(arrowYesNorthIdx, x, y, o)
}

func (b CircuitBoardRenderer) ArrowNo(x, y int, o Orientation) ImageToDraw {
	return b.arrow(arrowNoNorthIdx, x, y, o)
}

func (b CircuitBoardRenderer) imageByIndex(i int, x, y int) ImageToDraw {
	opts := ebiten.DrawImageOptions{}
	b.chips.Anchor(&opts.GeoM)
	opts.GeoM.Translate((float64(x)+0.5)*b.width, (float64(y)+0.5)*b.height)
	return ImageToDraw{
		Image:   b.chips.GetImage(i, 0),
		Options: &opts,
		Z:       chipZ,
	}
}

func (b CircuitBoardRenderer) arrow(baseIdx int, x, y int, o Orientation) ImageToDraw {
	// log.Printf("Arrow %d, x=%d, y=%d, o=%d", baseIdx, x, y, o)
	i := baseIdx + int(o)
	v := o.VelocityForward()
	opts := ebiten.DrawImageOptions{}
	b.chips.Anchor(&opts.GeoM)
	opts.GeoM.Translate(
		(float64(x)+0.5*(1+float64(v.Dx)))*b.width,
		(float64(y)+0.5*(1+float64(v.Dy)))*b.height,
	)
	return ImageToDraw{
		Image:   b.chips.GetImage(i, 0),
		Options: &opts,
		Z:       arrowZ,
	}
}

var chipType2imageIdx = map[ChipType]int{
	StartChip:       startIdx,
	ForwardChip:     forwardIdx,
	TurnLeftChip:    turnLeftIdx,
	TurnRightChip:   turnRightIdx,
	PaintRedChip:    paintRedIdx,
	PaintYellowChip: paintYellowIdx,
	PaintBlueChip:   paintBlueIdx,
	IsWallAheadChip: isWallAheadIdx,
	IsFloorRedChip:  isFloorRedIdx,
	IsFloorBlueChip: isFloorBueIdx,
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
