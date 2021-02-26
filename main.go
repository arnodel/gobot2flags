package main

import (
	"flag"
	_ "image/png"
	"io/ioutil"
	"log"

	"github.com/arnodel/gobot2flags/resources"
	"github.com/hajimehoshi/ebiten/v2"
)

const debug = false

const (
	frameWidth  = 32
	frameHeight = 32

	wallWidth  = 6
	wallHeight = 7
)

func main() {
	var levelFile string
	flag.StringVar(&levelFile, "level", "", "path to r2f level file")
	flag.Parse()
	var levelString = `
+--+--+--+--+
|RF|R |R  RF|
+  .  .  .  +
|Y  B> Y  B |
+--+--+  +  +
|BF Y  B |YF|
+--+--+--+--+	
`
	if levelFile != "" {
		levelBytes, err := ioutil.ReadFile(levelFile)
		if err != nil {
			log.Fatal("Could not open level file:", err)
		}
		levelString = string(levelBytes)
	}

	maze, err := MazeFromString(levelString)
	if err != nil {
		log.Fatalf("Could not create maze: %s", err)
	}
	// 	board, err := CircuitBoardFromString(`
	// |ST -> W? y> TL|
	// | ^    nv     v|
	// |..    MF <- ..|
	// | ^     v      |
	// |.. <- PR      |`)
	// 	if err != nil {
	// 		log.Fatalf("Could not create circuit board: %s", err)
	// 	}
	board := NewCircuitBoard(8, 8)
	chips := ChipRenderer{NewSprite(resources.GetImage("circuitboardtiles.png"), 32, 32, 16, 16)}
	boardRenderer := NewCircuitBoardRenderer(chips)
	mazeRenderer := &MazeRenderer{
		cellWidth:  frameWidth,
		cellHeight: frameHeight,
		wallWidth:  wallWidth,
		wallHeight: wallHeight,
		walls:      NewWalls(resources.GetImage("greywalls.png")),
		floors:     LoadFloors(resources.GetImage("floors.png")),
		robot:      NewSprite(resources.GetImage("robot.png"), frameWidth, frameHeight, 16, 16),
		flag:       NewSprite(resources.GetImage("greenflag.png"), frameWidth, frameHeight, 10, 28),
	}
	icons := Icons{NewSprite(resources.GetImage("icons.png"), 32, 32, 16, 16)}
	game := &Game{
		maze:            maze,
		mazeRenderer:    mazeRenderer,
		board:           board,
		boardRenderer:   &boardRenderer,
		controller:      &ManualController{},
		boardController: newBoardController(board, maze.Clone()),
		showBoard:       true,
		chipSelector: &boardTiles{
			selectedType: StartChip,
			icons:        icons,
		},
		gameControlSelector: &gameControlSelector{
			selectedControl: Rewind,
			icons:           icons,
		},
	}
	ebiten.SetWindowSize(1024, 768)
	ebiten.SetWindowTitle("Gobot 2 Flags")
	ebiten.SetWindowResizable(true)
	if err := ebiten.RunGame(game); err != nil {
		log.Fatal(err)
	}
}

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

type Icons struct {
	*Sprite
}

func (i Icons) Get(tp IconType) *ebiten.Image {
	return i.GetImage(int(tp), 0)
}
