package resources

import (
	"embed"
	"image"
	"log"

	"github.com/hajimehoshi/ebiten/v2"
)

//go:embed images
var resources embed.FS

// GetImage loads an image from the embedded filesystem and converts it to an
// ebiten image.
func GetImage(name string) *ebiten.Image {
	f, err := resources.Open("images/" + name)
	if err != nil {
		panic(err)
	}
	img, _, err := image.Decode(f)
	if err != nil {
		log.Fatal(err)
	}
	eimg := ebiten.NewImageFromImage(img)
	return eimg
}
