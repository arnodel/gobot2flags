package engine

import (
	"image"

	"github.com/hajimehoshi/ebiten/v2"
)

type Sprite struct {
	img                      *ebiten.Image
	width, height            int
	anchorX, anchorY         float64
	variantCount, frameCount int
}

func NewSprite(img *ebiten.Image, width, height int, anchorX, anchorY float64) *Sprite {
	iw, ih := img.Size()
	return &Sprite{
		img:          img,
		width:        width,
		height:       height,
		anchorX:      anchorX,
		anchorY:      anchorY,
		variantCount: ih / height,
		frameCount:   iw / width,
	}
}

func (s Sprite) GetImage(variant, frame int) *ebiten.Image {
	variant %= s.variantCount
	frame %= s.frameCount
	rect := image.Rect(frame*s.width, variant*s.height, (frame+1)*s.width, (variant+1)*s.height)
	return s.img.SubImage(rect).(*ebiten.Image)
}

func (s Sprite) ImageToDraw(variant, frame int) ImageToDraw {
	opts := ebiten.DrawImageOptions{}
	s.Anchor(&opts.GeoM)
	return ImageToDraw{
		Image:   s.GetImage(variant, frame),
		Options: &opts,
	}
}

func (s Sprite) Rotate(theta float64) ebiten.GeoM {
	g := ebiten.GeoM{}
	if theta != 0 {
		tx, ty := 0.5*float64(s.width), 0.5*float64(s.height)
		g.Translate(-tx, -ty)
		g.Rotate(theta)
		g.Translate(tx, ty)
	}
	return g
}

func (s Sprite) Anchor(g *ebiten.GeoM) {
	g.Translate(-s.anchorX, -s.anchorY)
}

func (s Sprite) Bounds() image.Rectangle {
	return image.Rect(0, 0, s.width, s.height).Sub(pt(s.anchorX, s.anchorY))
}
