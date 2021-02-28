package main

import (
	_ "image/png"
	"log"

	"github.com/arnodel/gobot2flags/engine"
	"github.com/arnodel/gobot2flags/play"
	"github.com/arnodel/gobot2flags/resources"
	"github.com/arnodel/gobot2flags/selectlevel"
	"github.com/hajimehoshi/ebiten/v2"
)

func main() {

	ebiten.SetWindowSize(1024, 768)
	ebiten.SetWindowTitle("Gobot 2 Flags")
	ebiten.SetWindowResizable(true)

	game := newGameController()
	if err := ebiten.RunGame(game); err != nil {
		log.Fatal(err)
	}
}

type gameController struct {
	levels     []string
	selectView engine.View
	playViews  map[string]engine.View
	engine.Game
}

func newGameController() *gameController {
	c := &gameController{
		levels:    resources.GetLevelList(),
		playViews: map[string]engine.View{},
	}
	c.setSelectView()
	return c
}

func (c *gameController) selectLevel(i int) {
	levelName := c.levels[i]
	level, err := resources.GetLevel(levelName)
	if err != nil {
		log.Println(err)
		return
	}
	playView := c.playViews[levelName]
	if playView == nil {
		playView = play.NewView(level, c.setSelectView)
		c.playViews[levelName] = playView
	}
	c.SetView(playView)
}

func (c *gameController) setSelectView() {
	if c.selectView == nil {
		c.selectView = selectlevel.NewView(c.levels, c.selectLevel)
	}
	c.SetView(c.selectView)
}
