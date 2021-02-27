package resources

import (
	"embed"
	"image"
	_ "image/png" // This is so that png type is registered with the image package and image.Decode() works
	"io/ioutil"

	"github.com/hajimehoshi/ebiten/v2"
	"golang.org/x/image/font/opentype"
)

//go:embed images fonts
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
		panic(err)
	}
	eimg := ebiten.NewImageFromImage(img)
	return eimg
}

func GetFont(name string) *opentype.Font {
	f, err := resources.Open("fonts/" + name)
	if err != nil {
		panic(err)
	}
	data, err := ioutil.ReadAll(f)
	if err != nil {
		panic(err)
	}
	font, err := opentype.Parse(data)
	if err != nil {
		panic(err)
	}
	return font
}
