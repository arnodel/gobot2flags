package main

import (
	"log"

	"github.com/hajimehoshi/ebiten/v2"
)

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

type BoardController struct {
	board    *CircuitBoard
	maze     *Maze
	robot    *Robot
	boardPos Position
}

func newBoardController(board *CircuitBoard, maze *Maze) *BoardController {
	return &BoardController{
		board:    board,
		maze:     maze,
		robot:    maze.robot,
		boardPos: board.startPos,
	}
}

func (c *BoardController) NextCommand() Command {
	c.board.ClearActiveChips()
	var (
		pos        = c.robot.Position
		wallAhead  = c.maze.HasWallAt(pos.X, pos.Y, c.robot.Orientation)
		floorColor = c.maze.CellAt(pos.X, pos.Y).Color()
	)
	for {
		var (
			chip            = c.board.ChipAt(c.boardPos.X, c.boardPos.Y)
			com, arrowType  = chip.Command(floorColor, wallAhead)
			nextChipDir, ok = chip.Arrow(arrowType)
		)
		c.board.ActivateChip(c.boardPos.X, c.boardPos.Y, nextChipDir)
		if ok {
			c.boardPos = c.boardPos.Move(nextChipDir.VelocityForward())
		} else {
			c.boardPos = c.board.startPos
		}
		log.Printf("board -> %v, com: %v", c.boardPos, com)
		if com != NoCommand || c.board.ChipAt(c.boardPos.X, c.boardPos.Y).IsActive() {
			return com
		}
	}
}
