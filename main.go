package main

import (
	"bytes"
	"image"
	_ "image/png"
	"log"

	"github.com/arnodel/gobots2flags/resources"
	"github.com/hajimehoshi/ebiten"
)

const (
	screenWidth  = 320
	screenHeight = 240

	frameOX     = 0
	frameOY     = 0
	frameWidth  = 32
	frameHeight = 32
	frameNum    = 6

	wallWidth  = 6
	wallHeight = 7
)

type Game struct {
	outsideWidth, outsideHeight int
	count                       int
	step                        int
	mazeRenderer                *MazeRenderer
	maze                        *Maze
	controller                  *ManualController
}

func (g *Game) Update(screen *ebiten.Image) error {
	g.count++
	if g.maze.robot.CurrentCommand != NoCommand {
		g.step = (g.step + 1) % 60
	}
	if g.step == 0 {
		g.controller.UpdateNextCommand()
		g.maze.AdvanceRobot(g.controller.NextCommand())
	}
	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	const scale = 2
	baseTr := ebiten.GeoM{}
	baseTr.Translate(-float64(g.maze.width*frameWidth)/2, -float64(g.maze.height*frameHeight)/2)
	baseTr.Scale(scale, scale)
	baseTr.Translate(float64(g.outsideWidth)/2, float64(g.outsideHeight)/2)

	canvas := &transformCanvas{
		target:   screen,
		baseGeoM: baseTr,
	}
	g.maze.Draw(canvas, g.mazeRenderer, float64(g.step)/60, g.count/60)
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	g.outsideWidth = outsideWidth
	g.outsideHeight = outsideHeight
	return outsideWidth, outsideHeight
}

func subImage(img *ebiten.Image, x, y int) *ebiten.Image {
	return img.SubImage(image.Rect(x*frameWidth, y*frameHeight, (x+1)*frameWidth, (y+1)*frameHeight)).(*ebiten.Image)
}

type Walls struct {
	Horizontal, Vertical, Corner *ebiten.Image
}

func NewWalls(img *ebiten.Image) Walls {
	return Walls{
		Horizontal: img.SubImage(image.Rect(0, frameHeight, frameWidth, frameHeight+wallHeight)).(*ebiten.Image),
		Vertical:   img.SubImage(image.Rect(0, 2*frameHeight, wallWidth, 3*frameHeight)).(*ebiten.Image),
		Corner:     img.SubImage(image.Rect(0, 0, wallWidth, wallHeight)).(*ebiten.Image),
	}
}

type Floors [4]*ebiten.Image

func LoadFloors(img *ebiten.Image) Floors {
	return Floors{
		subImage(img, 0, 0),
		subImage(img, 1, 0),
		subImage(img, 2, 0),
		subImage(img, 3, 0),
	}
}

func (f Floors) GetImage(c Color) *ebiten.Image {
	return f[c]
}

func getImage(b []byte) *ebiten.Image {
	img, _, err := image.Decode(bytes.NewReader(b))
	if err != nil {
		log.Fatal(err)
	}
	eimg, _ := ebiten.NewImageFromImage(img, ebiten.FilterDefault)
	return eimg
}

func main() {
	maze, err := MazeFromString(`
+--+--+--+--+
|RF|R |R  R |
+  .  .  .  +
|Y  B> Y  B |
+--+--+  +  +
|B  Y  B |YF|
+--+--+--+--+
`)
	if err != nil {
		log.Fatalf("Could not create maze: %s", err)
	}
	mazeRenderer := &MazeRenderer{
		cellWidth:  frameWidth,
		cellHeight: frameHeight,
		wallWidth:  wallWidth,
		wallHeight: wallHeight,
		walls:      NewWalls(getImage(resources.GreyWallsPng)),
		floors:     LoadFloors(getImage(resources.FloorsPng)),
		robot:      NewSprite(getImage(resources.RobotPng), frameWidth, frameHeight, 16, 16),
		flag:       NewSprite(getImage(resources.GreenFlagPng), frameWidth, frameHeight, 10, 28),
	}
	game := &Game{
		maze:         maze,
		mazeRenderer: mazeRenderer,
		controller:   &ManualController{},
	}
	ebiten.SetWindowSize(screenWidth*2, screenHeight*2)
	ebiten.SetWindowTitle("Gobot 2 Flags")
	if err := ebiten.RunGame(game); err != nil {
		log.Fatal(err)
	}
}
