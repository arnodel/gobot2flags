package main

import (
	"github.com/arnodel/gobot2flags/resources"
	"golang.org/x/image/font"
	"golang.org/x/image/font/opentype"
)

var msgFont font.Face

func init() {
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
