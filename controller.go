package main

import "github.com/hajimehoshi/ebiten/v2"

type Command int

const (
	NoCommand Command = iota
	TurnLeft
	TurnRight
	MoveForward
	PaintRed
	PaintBlue
	PaintYellow
)

type ManualController struct {
	nextCommand  Command
	spacePressed bool
	switchView   bool
}

func (c *ManualController) NextCommand() (com Command) {
	com = c.nextCommand
	c.nextCommand = NoCommand
	return
}

func (c *ManualController) GetSwitch() bool {
	s := c.switchView
	c.switchView = false
	return s
}

func (c *ManualController) UpdateNextCommand() Command {
	spacePressed := ebiten.IsKeyPressed(ebiten.KeySpace)
	if spacePressed && !c.spacePressed {
		c.switchView = true
	}
	c.spacePressed = spacePressed
	switch {
	case ebiten.IsKeyPressed(ebiten.KeyLeft):
		c.nextCommand = TurnLeft
	case ebiten.IsKeyPressed(ebiten.KeyUp):
		c.nextCommand = MoveForward
	case ebiten.IsKeyPressed(ebiten.KeyRight):
		c.nextCommand = TurnRight
	}
	return c.nextCommand
}
