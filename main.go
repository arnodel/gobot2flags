package main

import (
	"flag"
	_ "image/png"
	"io/ioutil"
	"log"

	"github.com/arnodel/gobot2flags/engine"
	"github.com/arnodel/gobot2flags/model"
	"github.com/arnodel/gobot2flags/resources"
	"github.com/arnodel/gobot2flags/sprites"
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

	maze, err := model.MazeFromString(levelString)
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
	board := model.NewCircuitBoard(8, 8)
	chips := ChipRenderer{engine.NewSprite(resources.GetImage("circuitboardtiles.png"), 32, 32, 16, 16)}
	boardRenderer := NewCircuitBoardRenderer(chips)
	mazeRenderer := &MazeRenderer{
		cellWidth:  frameWidth,
		cellHeight: frameHeight,
		wallWidth:  wallWidth,
		wallHeight: wallHeight,
		walls:      sprites.GreyWalls,
		floors:     sprites.PlainFloors,
		robot:      sprites.Robot,
		flag:       sprites.Flag,
	}
	game := &Game{
		maze:            maze,
		mazeRenderer:    mazeRenderer,
		board:           board,
		boardRenderer:   &boardRenderer,
		boardController: model.NewBoardController(board, maze.Clone()),
		showBoard:       true,
		chipSelector: &boardTiles{
			selectedType: model.StartChip,
			icons:        sprites.PlainIcons,
		},
		gameControlSelector: &gameControlSelector{
			selectedControl: Rewind,
			icons:           sprites.PlainIcons,
		},
		pointer: &engine.PointerTracker{},
	}
	ebiten.SetWindowSize(1024, 768)
	ebiten.SetWindowTitle("Gobot 2 Flags")
	ebiten.SetWindowResizable(true)
	if err := ebiten.RunGame(game); err != nil {
		log.Fatal(err)
	}
}
