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

func (c Command) String() string {
	switch c {
	case NoCommand:
		return "no command"
	case TurnLeft:
		return "turn left"
	case TurnRight:
		return "turn right"
	case MoveForward:
		return "move forward"
	case PaintRed:
		return "paint floor red"
	case PaintBlue:
		return "paint floor blue"
	case PaintYellow:
		return "paint floor yellow"
	default:
		return "unknown"
	}
}

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

func (c *BoardController) GameWon() bool {
	return c.maze.FlagsRemaining() == 0
}

func (c *BoardController) Advance() {
	if c.GameWon() {
		log.Printf("Game won!")
		c.board.ClearActiveChips()
		return
	}
	c.maze.AdvanceRobot(c.NextCommand())
}

func (c *BoardController) NextCommand() Command {
	c.board.ClearActiveChips()
	var (
		pos        = c.robot.Position
		wallAhead  = c.maze.HasWallAt(pos.X, pos.Y, c.robot.Orientation)
		floorColor = c.maze.CellAt(pos.X, pos.Y).Color()
	)
	log.Printf("Robot at %s, facing %s", c.robot.Position, c.robot.Orientation)
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
		log.Printf("Board -> %s, com: %s", c.boardPos, com)
		if com != NoCommand || c.board.ChipAt(c.boardPos.X, c.boardPos.Y).IsActive() {
			return com
		}
	}
}
