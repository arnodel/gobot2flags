package main

import (
	"errors"
	"fmt"
	"log"
	"strings"
)

const (
	TFPos = iota
	LFPos
	CFPos
	FFPos
	C1Pos
	C2Pos
)

type Cell byte

const (
	TF Cell = 1 << iota // North Wall
	LF                  // West Wall
	CF                  // Corner Wall (Northwest)
	FF                  // Flag
	C1                  // Floor Color 1
	C2                  // Floor Color 2
	KF                  // Captured flag
)

type Color byte

const (
	NoColor Color = iota
	Red
	Yellow
	Blue
)

func (c Color) ToCell() Cell {
	return Cell(c << C1Pos)
}

func (c Cell) NorthWall() bool {
	return (c & TF) != 0
}

func (c Cell) WestWall() bool {
	return (c & LF) != 0
}

func (c Cell) CornerWall() bool {
	return (c & CF) != 0
}

func (c Cell) Color() Color {
	return Color((c & (C1 | C2)) >> C1Pos)
}

func (c Cell) Flag() bool {
	return (c & FF) != 0
}

func (c Cell) Captured() bool {
	return (c & KF) != 0
}

type Maze struct {
	width, height int
	cells         []Cell
	robot         *Robot
}

func NewMaze(width, height int) *Maze {
	return &Maze{
		width:  width,
		height: height,
		cells:  make([]Cell, width*height),
	}
}

func (m *Maze) cellIndex(x, y int) int {
	if x < 0 {
		x = x%m.width + m.width
	} else if x >= m.width {
		x = x % m.width
	}
	if y < 0 {
		y = y%m.height + m.height
	} else if y >= m.height {
		y = y % m.height
	}
	return x + y*m.width
}

func (m *Maze) UpdateCellAt(x, y int, c Cell) {
	m.cells[m.cellIndex(x, y)] |= c
}

func (m *Maze) CellAt(x, y int) Cell {
	return m.cells[m.cellIndex(x, y)]
}

func (m *Maze) HasWallAt(x, y int, o Orientation) bool {
	switch o {
	case North:
		return m.CellAt(x, y).NorthWall()
	case West:
		return m.CellAt(x, y).WestWall()
	case South:
		return m.CellAt(x, y+1).NorthWall()
	case East:
		return m.CellAt(x+1, y).WestWall()
	default:
		return false // Or panic?
	}
}

func wrongCharErr(i, j int, allowed string) error {
	return fmt.Errorf("wrong char line %d col %d: one of [%s] allowed", i+1, j+1, allowed)
}

func MazeFromString(s string) (*Maze, error) {
	s = strings.TrimSpace(s)
	rows := strings.Split(s, "\n")
	if len(rows)%2 != 1 {
		return nil, errors.New("need odd number of lines")
	}
	height := (len(rows) - 1) / 2
	if height == 0 {
		return nil, errors.New("need at least 1 row")
	}
	lr0 := len(rows[0])
	width := (lr0 - 1) / 3
	if width*3+1 != lr0 {
		return nil, errors.New("wrong length for line 1")
	}
	for i, row := range rows {
		if len(row) != lr0 {
			return nil, fmt.Errorf("wrong length for line %d", i+1)
		}
	}
	maze := NewMaze(width, height)
	for i, row := range rows {
		y := i / 2
		if y == height {
			// Ignore last row (should validate)
			continue
		}
		if i%2 == 0 {
			// Parsing a horizontal wall row
			for j, c := range row {
				x := j / 3
				if x == width {
					if c != '+' {
						return nil, wrongCharErr(i, j, "+")
					}
					continue
				}
				switch j % 3 {
				case 0:
					// Corner
					switch c {
					case '+':
						maze.UpdateCellAt(x, y, CF)
					case '.':
						// No corner
					default:
						return nil, wrongCharErr(i, j, "+.")
					}
				case 1:
					// Horizontal wall
					switch c {
					case '-':
						maze.UpdateCellAt(x, y, TF)
					case ' ':
						// No wall
					default:
						return nil, wrongCharErr(i, j, "- ")
					}
				case 2:
					// Check agrees with previous one
					switch c {
					case '-':
						if !maze.CellAt(x, y).NorthWall() {
							return nil, wrongCharErr(i, j, " ")
						}
					case ' ':
						if maze.CellAt(x, y).NorthWall() {
							return nil, wrongCharErr(i, j, "-")
						}
					default:
						return nil, wrongCharErr(i, j, "- ")
					}
				}
			}
		} else {
			// Parsing a floor row
			for j, c := range row {
				x := j / 3
				if x == width {
					if c != '|' {
						return nil, wrongCharErr(i, j, "|")
					}
					continue
				}
				switch j % 3 {
				case 0:
					// Vertical wall
					switch c {
					case '|':
						maze.UpdateCellAt(x, y, LF)
					case ' ':
						// No wall
					default:
						return nil, wrongCharErr(i, j, "| ")
					}
				case 1:
					// Floor
					switch c {
					case 'R':
						maze.UpdateCellAt(x, y, Red.ToCell())
					case 'Y':
						maze.UpdateCellAt(x, y, Yellow.ToCell())
					case 'B':
						maze.UpdateCellAt(x, y, Blue.ToCell())
					case ' ':
						// No color
					default:
						return nil, wrongCharErr(i, j, "RYB ")
					}
				case 2:
					// Flag or robot
					switch c {
					case 'F':
						maze.UpdateCellAt(x, y, FF)
					case '>', '<', '^', 'v':
						if maze.robot != nil {
							return nil, fmt.Errorf("Only one robot allowed (second defined line %d, col %d)", i, j)
						}
						maze.robot = &Robot{
							Position: Position{
								X: x,
								Y: y,
							},
							Orientation: rune2Orientation[c],
						}
					case ' ':
						// Nothing
					default:
						return nil, wrongCharErr(i, j, "F ")
					}
				}
			}
		}
	}
	return maze, nil
}

func (m *Maze) Draw(c Canvas, r *MazeRenderer, t float64, frame int) {
	h := m.height
	w := m.width

	// Draw the floors first as they are under everything
	for y := 0; y < h; y++ {
		for x := 0; x < w; x++ {
			r.DrawFloor(c, x, y, m.CellAt(x, y).Color())
		}
	}

	robot := m.robot
	stack := ImageStack{}

	// Draw the walls and objects
	for y := 0; y <= h; y++ {
		for x := 0; x <= w; x++ {
			cell := m.CellAt(x%w, y%h)
			if cell.CornerWall() {
				stack.Add(r.CornerWall(x, y))
			}
			if y < h && cell.WestWall() {
				stack.Add(r.WestWall(x, y))
			}
			if x < w && cell.NorthWall() {
				stack.Add(r.NorthWall(x, y))
			}
			if x < w && y < h && cell.Flag() {
				stack.Add(r.Flag(x, y, frame, cell.Captured()))
			}
		}
		// Draw the robot now (so it's placed correctly w.r.t. walls)
		if robot != nil && robot.Y == y {
			stack.Add(r.Robot(robot, t, frame))
		}
		stack.Draw(c)
		stack.Empty() // Reuse the memory
	}
}

func (m *Maze) AdvanceRobot(com Command) bool {
	robot := m.robot.Advance()
	cell := m.CellAt(robot.X, robot.Y)
	if cell.Flag() && !cell.Captured() {
		log.Printf("capture %d %d", robot.X, robot.Y)
		m.UpdateCellAt(robot.X, robot.Y, KF)
	}
	next := robot.ApplyCommand(com)
	crash := next.IsMovingForward() && m.HasWallAt(next.X, next.Y, next.Orientation)
	if crash {
		next = robot.ApplyCommand(NoCommand)
	}
	*m.robot = next
	return !crash
}

var rune2Orientation = map[rune]Orientation{
	'>': East,
	'<': West,
	'^': North,
	'v': South,
}
