package main

import "github.com/hajimehoshi/ebiten"

type CircuitBoardRenderer struct {
	chips *Sprite
}

func NewBoardChipImages(img *ebiten.Image) CircuitBoardRenderer {
	return CircuitBoardRenderer{chips: NewSprite(img, 32, 32, 16, 16)}
}

func (b CircuitBoardRenderer) Chip(c ChipType, x, y int) (toDraw ImageToDraw, ok bool) {
	var i int
	i, ok = chipType2imageIdx[c]
	if ok {
		toDraw.Image = b.chips.GetImage(i, 0)
		b.chips.Anchor(&toDraw.Options.GeoM)
		toDraw.Z = chipZ
	}
	return
}

func (b CircuitBoardRenderer) Background() (toDraw ImageToDraw) {
	toDraw.Image = b.chips.GetImage(backgroundIdx, 0)
	b.chips.Anchor(&toDraw.Options.GeoM)
	toDraw.Z = bgZ
	return
}

var chipType2imageIdx = map[ChipType]int{
	StartChip:       startIdx,
	ForwardChip:     forwardIdx,
	TurnLeftChip:    turnLeftIdx,
	TurnRightChip:   turnRightIdx,
	PaintRedChip:    paintRedIdx,
	PaintYellowChip: paintYellowIdx,
	PaintBlueChip:   painBlueIdx,
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
	painBlueIdx
)

const (
	bgZ float64 = iota
	arrowZ
	chipZ
)
