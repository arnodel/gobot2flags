package resources

import (
	"embed"
	"image"
	_ "image/png" // This is so that png type is registered with the image package and image.Decode() works
	"io/ioutil"
	"strings"

	"github.com/arnodel/gobot2flags/model"
	"github.com/hajimehoshi/ebiten/v2"
	"golang.org/x/image/font/opentype"
)

//go:embed images fonts levels
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

func GetLevelList() []string {
	entries, err := resources.ReadDir("levels")
	if err != nil {
		panic(err)
	}
	var levels []string
	for _, entry := range entries {
		if entry.Type().IsRegular() && strings.HasSuffix(entry.Name(), ".r2f") {
			levels = append(levels, strings.TrimSuffix(entry.Name(), ".r2f"))
		}
	}
	return levels
}

func GetLevel(name string) (*model.Level, error) {
	f, err := resources.Open("levels/" + name + ".r2f")
	if err != nil {
		return nil, err
	}
	data, err := ioutil.ReadAll(f)
	if err != nil {
		return nil, err
	}
	return model.LevelFromString(name, string(data))
}
