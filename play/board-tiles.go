package play

import (
	"image"

	"github.com/arnodel/gobot2flags/engine"
	"github.com/arnodel/gobot2flags/model"
	"github.com/arnodel/gobot2flags/sprites"
)

var chipTypes = []model.ChipType{
	model.StartChip,
	model.ForwardChip,
	model.TurnLeftChip,
	model.TurnRightChip,
	model.PaintRedChip,
	model.PaintYellowChip,
	model.PaintBlueChip,
	model.IsWallAheadChip,
	model.IsFloorRedChip,
	model.IsFloorYellowChip,
	model.IsFloorBlueChip,
}

var arrowTypes = []model.ArrowType{
	model.ArrowYes,
	model.ArrowNo,
}

var boardIcons = []sprites.IconType{
	sprites.EraserIcon,
	sprites.TrashCanIcon,
}

var boardTilesImages []engine.ImageToDraw

type boardTiles struct {
	selectedType      model.ChipType
	selectedArrowType model.ArrowType
	selectedIcon      sprites.IconType
	indexSelector     engine.Selector
	selectedIndex     int
	grid              engine.Grid
}

func (b *boardTiles) Bounds() image.Rectangle {
	return b.grid.Bounds()
}

func (b *boardTiles) Draw(c engine.Canvas, chips ChipRenderer) {
	for i, img := range boardTilesImages {
		opts := *img.Options
		var scale float64
		p := b.indexSelector.SelectingProportion()
		switch {
		case b.indexSelector.IsSelecting(i):
			scale = (1 + p) / 2
		case b.selectedIndex == i:
			scale = 1 - p/2
		default:
			scale = 0.5
		}
		opts.GeoM.Scale(scale, scale)
		opts.GeoM.Translate(b.grid.CellCenter(i, 0))
		c.DrawImage(img.Image, &opts)
	}
}

func (b *boardTiles) Update(p engine.PointerStatus) {
	nControls := len(boardTilesImages)
	b.grid = engine.Grid{
		CellWidth:  24,
		CellHeight: 32,
		Rows:       2,
		Columns:    (nControls + 1) / 2,
	}

	idx := b.grid.CellIndex(p.CurrentCoords())
	if b.indexSelector.Update(idx, p.Status()) == engine.Select {
		b.selectIndex(idx)
	}
}

func (b *boardTiles) selectIndex(idx int) {
	b.selectedIndex = idx

	selectedArrowType := model.NoArrow
	selectedType := model.NoChip
	selectedIcon := sprites.NoIcon

	switch {
	case idx < 0:
		break
	case idx < len(chipTypes):
		selectedType = chipTypes[idx]
	case idx < len(chipTypes)+len(arrowTypes):
		selectedArrowType = arrowTypes[idx-len(chipTypes)]
	case idx < len(chipTypes)+len(arrowTypes)+len(boardIcons):
		selectedIcon = boardIcons[idx-len(chipTypes)-len(arrowTypes)]
	}

	b.selectedArrowType = selectedArrowType
	b.selectedType = selectedType
	b.selectedIcon = selectedIcon
}
