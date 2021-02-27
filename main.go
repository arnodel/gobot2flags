package main

import (
	"flag"
	_ "image/png"
	"io/ioutil"
	"log"

	"github.com/arnodel/gobot2flags/game"
	"github.com/arnodel/gobot2flags/model"
	"github.com/hajimehoshi/ebiten/v2"
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
	theGame := game.New(maze, board)

	ebiten.SetWindowSize(1024, 768)
	ebiten.SetWindowTitle("Gobot 2 Flags")
	ebiten.SetWindowResizable(true)
	if err := ebiten.RunGame(theGame); err != nil {
		log.Fatal(err)
	}
}
