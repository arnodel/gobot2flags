package sprites

import (
	"image"

	"github.com/arnodel/gobot2flags/engine"
	"github.com/arnodel/gobot2flags/model"
	"github.com/arnodel/gobot2flags/resources"
	"github.com/hajimehoshi/ebiten/v2"
)

const (
	FrameWidth  = 32
	FrameHeight = 32

	WallWidth  = 6
	WallHeight = 7
)

type IconType int

const (
	PlayIcon IconType = iota
	StepIcon
	FastForwardIcon
	RewindIcon
	PauseIcon
	TrashCanIcon
	EraserIcon
)

const NoIcon IconType = -1

var (
	GreyWalls         Walls
	PlainFloors       Floors
	Robot             *engine.Sprite
	Flag              *engine.Sprite
	CircuitBoardTiles *engine.Sprite
	PlainIcons        Icons
)

func init() {
	GreyWalls = newWalls(resources.GetImage("greywalls.png"))
	PlainFloors = loadFloors(resources.GetImage("floors.png"))
	Robot = engine.NewSprite(resources.GetImage("robot.png"), FrameWidth, FrameHeight, FrameWidth/2, FrameHeight/2)
	Flag = engine.NewSprite(resources.GetImage("greenflag.png"), FrameWidth, FrameHeight, 10, 28)
	PlainIcons = Icons{engine.NewSprite(resources.GetImage("icons.png"), 32, 32, 16, 16)}
	CircuitBoardTiles = engine.NewSprite(resources.GetImage("circuitboardtiles.png"), 32, 32, 16, 16)
}

type Walls struct {
	Horizontal, Vertical, Corner *ebiten.Image
}

func newWalls(img *ebiten.Image) Walls {
	return Walls{
		Horizontal: img.SubImage(image.Rect(0, FrameHeight, FrameWidth, FrameHeight+WallHeight)).(*ebiten.Image),
		Vertical:   img.SubImage(image.Rect(0, 2*FrameHeight, WallWidth, 3*FrameHeight)).(*ebiten.Image),
		Corner:     img.SubImage(image.Rect(0, 0, WallWidth, WallHeight)).(*ebiten.Image),
	}
}

type Floors [4]*ebiten.Image

func (f Floors) GetImage(c model.Color) *ebiten.Image {
	return f[c]
}

// TODO: turn to a sprite
func loadFloors(img *ebiten.Image) Floors {
	return Floors{
		subImage(img, 0, 0),
		subImage(img, 1, 0),
		subImage(img, 2, 0),
		subImage(img, 3, 0),
	}
}

type Icons struct {
	*engine.Sprite
}

func (i Icons) Get(tp IconType) *ebiten.Image {
	return i.GetImage(int(tp), 0)
}

func subImage(img *ebiten.Image, x, y int) *ebiten.Image {
	return img.SubImage(image.Rect(x*FrameWidth, y*FrameHeight, (x+1)*FrameWidth, (y+1)*FrameHeight)).(*ebiten.Image)
}
