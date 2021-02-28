package main

import (
	"flag"
	_ "image/png"
	"log"

	"github.com/arnodel/gobot2flags/engine"
	"github.com/arnodel/gobot2flags/play"
	"github.com/arnodel/gobot2flags/resources"
	"github.com/arnodel/gobot2flags/selectlevel"
	"github.com/hajimehoshi/ebiten/v2"
)

func main() {
	var levelFile string
	flag.StringVar(&levelFile, "level", "", "path to r2f level file")
	flag.Parse()
	// 	var levelString = `
	// +--+--+--+--+
	// |RF|R |R  RF|
	// +  .  .  .  +
	// |Y  B> Y  B |
	// +--+--+  +  +
	// |BF Y  B |YF|
	// +--+--+--+--+
	// `
	// if levelFile != "" {
	// 	levelBytes, err := ioutil.ReadFile(levelFile)
	// 	if err != nil {
	// 		log.Fatal("Could not open level file:", err)
	// 	}
	// 	levelString = string(levelBytes)
	// }

	// maze, err := model.MazeFromString(levelString)
	// if err != nil {
	// 	log.Fatalf("Could not create maze: %s", err)
	// }
	// 	board, err := CircuitBoardFromString(`
	// |ST -> W? y> TL|
	// | ^    nv     v|
	// |..    MF <- ..|
	// | ^     v      |
	// |.. <- PR      |`)
	// 	if err != nil {
	// 		log.Fatalf("Could not create circuit board: %s", err)
	// 	}
	levels := resources.GetLevelList()
	playViews := map[string]*play.View{}
	game := engine.NewGame(nil)
	var selectView engine.View
	var goToSelect = func() {
		game.SetView(selectView)
	}
	var selectLevel = func(i int) {
		level, err := resources.GetLevel(levels[i])
		if err != nil {
			log.Println(err)
			return
		}
		playView := playViews[levels[i]]
		if playView == nil {
			playView = play.New(level, goToSelect)
			playViews[levels[i]] = playView
		}
		game.SetView(playView)
	}
	selectView = selectlevel.NewView(levels, selectLevel)
	game.SetView(selectView)
	ebiten.SetWindowSize(1024, 768)
	ebiten.SetWindowTitle("Gobot 2 Flags")
	ebiten.SetWindowResizable(true)
	if err := ebiten.RunGame(game); err != nil {
		log.Fatal(err)
	}
}
