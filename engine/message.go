package engine

import (
	"image/color"

	"github.com/arnodel/gobot2flags/resources"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text"
	"golang.org/x/image/font"
	"golang.org/x/image/font/opentype"
)

var msgFont font.Face

func init() {
	// TODO: move this out of engine
	tt := resources.GetFont("c64.otf")
	const dpi = 72
	var err error
	msgFont, err = opentype.NewFace(tt, &opentype.FaceOptions{
		Size:    16,
		DPI:     dpi,
		Hinting: font.HintingFull,
	})
	if err != nil {
		panic(err)
	}
}

func DrawText(dest *ebiten.Image, txt string, x, y int, clr color.Color) {
	text.Draw(dest, txt, msgFont, x, y, clr)
}
