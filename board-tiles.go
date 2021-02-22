package main

import (
	"image"

	"github.com/hajimehoshi/ebiten/v2"
)

var chipTypes1 = []ChipType{
	StartChip,
	ForwardChip,
	TurnLeftChip,
	TurnRightChip,
	PaintRedChip,
	PaintYellowChip,
	PaintBlueChip,
}

var chipTypes2 = []ChipType{
	IsWallAheadChip,
	IsFloorRedChip,
	IsFloorYellowChip,
	IsFloorBlueChip,
}

var arrowTypes = []ArrowType{
	ArrowYes,
	ArrowNo,
}

var boardIcons = []IconType{EraserIcon, TrashCanIcon}

type boardTiles struct {
	selectedType      ChipType
	selectedArrowType ArrowType
	selectedIcon      IconType
	icons             Icons
}

func (b *boardTiles) Bounds() image.Rectangle {
	l1 := len(chipTypes1)
	l2 := len(chipTypes2) + len(arrowTypes) + len(boardIcons)
	maxLen := l1
	if l2 > l1 {
		maxLen = l2
	}
	return image.Rect(0, 0, maxLen*24, 64)
}

func (b *boardTiles) Draw(c Canvas, chips ChipRenderer) {
	for i, chipType := range chipTypes1 {
		img := chips.GetChipImage(chipType)
		var opts ebiten.DrawImageOptions
		chips.Anchor(&opts.GeoM)
		if b.selectedType != chipType {
			opts.GeoM.Scale(0.5, 0.5)
		}
		opts.GeoM.Translate((float64(i)+0.5)*24, 16)
		c.DrawImage(img, &opts)
	}
	for i, chipType := range chipTypes2 {
		img := chips.GetChipImage(chipType)
		var opts ebiten.DrawImageOptions
		chips.Anchor(&opts.GeoM)
		if b.selectedType != chipType {
			opts.GeoM.Scale(0.5, 0.5)
		}
		opts.GeoM.Translate((float64(i)+0.5)*24, 16+24)
		c.DrawImage(img, &opts)
	}
	start := len(chipTypes2)
	for i, arrowType := range arrowTypes {
		img := chips.GetArrowImage(arrowType)
		var opts ebiten.DrawImageOptions
		chips.Anchor(&opts.GeoM)
		if b.selectedArrowType != arrowType {
			opts.GeoM.Scale(0.5, 0.5)
		}
		opts.GeoM.Translate((float64(start+i)+0.5)*24, 16+24)
		c.DrawImage(img, &opts)
	}
	start += len(arrowTypes)
	for i, iconType := range boardIcons {
		img := b.icons.Get(iconType)
		var opts ebiten.DrawImageOptions
		chips.Anchor(&opts.GeoM)
		if b.selectedIcon != iconType {
			opts.GeoM.Scale(0.5, 0.5)
		}
		opts.GeoM.Translate((float64(start+i)+0.5)*24, 16+24)
		c.DrawImage(img, &opts)
	}
}

func (b *boardTiles) Click(x, y float64) {
	i := int(x / 24)
	j := int(y / 32)
	selectedArrowType := NoArrow
	selectedType := NoChip
	selectedIcon := NoIcon

	switch j {
	case 0:
		switch {
		case i < 0:
			return
		case i < len(chipTypes1):
			selectedType = chipTypes1[i]
		default:
			return
		}
	case 1:
		switch {
		case i < 0:
			return
		case i < len(chipTypes2):
			selectedType = chipTypes2[i]
		case i < len(chipTypes2)+len(arrowTypes):
			selectedArrowType = arrowTypes[i-len(chipTypes2)]
		case i < len(chipTypes2)+len(arrowTypes)+len(boardIcons):
			selectedIcon = boardIcons[i-len(chipTypes2)-len(arrowTypes)]
		default:
			return
		}
	default:
		return
	}
	b.selectedArrowType = selectedArrowType
	b.selectedType = selectedType
	b.selectedIcon = selectedIcon
}
