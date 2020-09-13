package main

import "github.com/hajimehoshi/ebiten"

type Command int

const (
	NoCommand Command = iota
	TurnLeft
	TurnRight
	MoveForward
)

type ManualController struct {
	nextCommand Command
}

func (c *ManualController) NextCommand() (com Command) {
	com = c.nextCommand
	c.nextCommand = NoCommand
	return
}

func (c *ManualController) UpdateNextCommand() Command {
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
