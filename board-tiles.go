package main

import (
	"image"

	"github.com/hajimehoshi/ebiten/v2"
)

var chipTypes = []ChipType{
	StartChip,
	ForwardChip,
	TurnLeftChip,
	TurnRightChip,
	PaintRedChip,
	PaintYellowChip,
	PaintBlueChip,
	IsWallAheadChip,
	IsFloorRedChip,
	IsFloorYellowChip,
	IsFloorBlueChip,
}

type boardTiles struct {
	selectedType      ChipType
	selectedArrowType ArrowType
}

var arrowTypes = []ArrowType{
	ArrowYes,
	ArrowNo,
}

func (b *boardTiles) Bounds() image.Rectangle {
	return image.Rect(0, 0, (len(arrowTypes)+len(chipTypes))*24, 32)
}

func (b *boardTiles) Draw(c Canvas, chips ChipRenderer) {
	for i, arrowType := range arrowTypes {
		img := chips.GetArrowImage(arrowType)
		var opts ebiten.DrawImageOptions
		chips.Anchor(&opts.GeoM)
		if b.selectedArrowType != arrowType {
			opts.GeoM.Scale(0.5, 0.5)
		}
		opts.GeoM.Translate((float64(i)+0.5)*24, 16)
		c.DrawImage(img, &opts)
	}
	for i, chipType := range chipTypes {
		img := chips.GetChipImage(chipType)
		var opts ebiten.DrawImageOptions
		chips.Anchor(&opts.GeoM)
		if b.selectedType != chipType {
			opts.GeoM.Scale(0.5, 0.5)
		}
		opts.GeoM.Translate((float64(i+len(arrowTypes))+0.5)*24, 16)
		c.DrawImage(img, &opts)
	}
}

func (b *boardTiles) Click(x, y float64) {
	idx := int(x / 24)
	if idx >= 0 && idx < len(arrowTypes) {
		selectedArrowType := arrowTypes[idx]
		if selectedArrowType == b.selectedArrowType {
			selectedArrowType = NoArrow
		}
		b.selectedArrowType = selectedArrowType
		return
	}
	idx -= 2
	if idx >= 0 && idx < len(chipTypes) {
		selectedType := chipTypes[idx]
		if selectedType == b.selectedType {
			selectedType = NoChip
		}
		b.selectedType = selectedType
	}
}
