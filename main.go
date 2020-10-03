package main

import (
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
	showBoard                   bool
	mazeRenderer                *MazeRenderer
	maze                        *Maze
	boardRenderer               *CircuitBoardRenderer
	board                       *CircuitBoard
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
	if g.showBoard {
		g.drawBoard(screen)
	} else {
		g.drawMaze(screen)
	}
}

func (g *Game) drawMaze(screen *ebiten.Image) {
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

func (g *Game) drawBoard(screen *ebiten.Image) {
	const scale = 2
	baseTr := ebiten.GeoM{}
	baseTr.Translate(-float64(g.board.width*frameWidth)/2, -float64(g.board.height*frameHeight)/2)
	baseTr.Scale(scale, scale)
	baseTr.Translate(float64(g.outsideWidth)/2, float64(g.outsideHeight)/2)

	canvas := &transformCanvas{
		target:   screen,
		baseGeoM: baseTr,
	}
	g.board.Draw(canvas, g.boardRenderer)
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	g.outsideWidth = outsideWidth
	g.outsideHeight = outsideHeight
	return outsideWidth, outsideHeight
}

func main() {
	maze, err := MazeFromString(`
+--+--+--+--+
|RF|R |R  RF|
+  .  .  .  +
|Y  B> Y  B |
+--+--+  +  +
|BF Y  B |YF|
+--+--+--+--+
`)
	if err != nil {
		log.Fatalf("Could not create maze: %s", err)
	}
	board, err := CircuitBoardFromString(`
|ST -> W? y> TL|
| ^    nv     v|
|..    MF <- ..|
| ^     v      |
|.. <- PR      |`)
	if err != nil {
		log.Fatalf("Could not create circuit board: %s", err)
	}
	boardRenderer := NewBoardChipImages(getImage(resources.CircuitBoarTiles))
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
		maze:          maze,
		mazeRenderer:  mazeRenderer,
		board:         board,
		boardRenderer: &boardRenderer,
		controller:    &ManualController{},
		showBoard:     true,
	}
	ebiten.SetWindowSize(screenWidth*2, screenHeight*2)
	ebiten.SetWindowTitle("Gobot 2 Flags")
	if err := ebiten.RunGame(game); err != nil {
		log.Fatal(err)
	}
}
