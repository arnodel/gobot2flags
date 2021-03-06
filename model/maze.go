package model

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

func (c Color) String() string {
	switch c {
	case NoColor:
		return "no color"
	case Red:
		return "red"
	case Yellow:
		return "yellow"
	case Blue:
		return "blue"
	default:
		return "unknown color"
	}
}

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

func (c Cell) SetColor(col Color) Cell {
	return (c & ^(C1 | C2)) | col.ToCell()
}

func (c Cell) Capture() Cell {
	return c | KF
}

type Maze struct {
	width, height   int
	cells           []Cell
	robot           *Robot
	flags, captured int
}

func NewMaze(width, height int) *Maze {
	return &Maze{
		width:  width,
		height: height,
		cells:  make([]Cell, width*height),
	}
}

func (m *Maze) Clone() *Maze {
	clone := NewMaze(m.width, m.height)
	copy(clone.cells, m.cells)
	robot := *m.robot
	clone.robot = &robot
	clone.flags = m.flags
	clone.captured = m.captured
	return clone
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

// TODO: remove this no-good function.
func (m *Maze) UpdateCellAt(x, y int, c Cell) {
	p := &m.cells[m.cellIndex(x, y)]
	c0c := *p &^ c
	cc0 := c &^ *p
	if c0c.Flag() {
		m.flags--
	} else if cc0.Flag() {
		m.flags++
	}
	if c0c.Captured() {
		m.captured--
	} else if cc0.Captured() {
		m.captured++
	}
	if c.Color() != NoColor {
		*p &= ^(C1 | C2)
	}
	*p |= c
}

func (m *Maze) CaptureFlag(x, y int) {
	p := &m.cells[m.cellIndex(x, y)]
	if !p.Flag() || p.Captured() {
		return
	}
	m.captured++
	*p = p.Capture()
}

func (m *Maze) PaintCell(x, y int, col Color) {
	p := &m.cells[m.cellIndex(x, y)]
	*p = p.SetColor(col)
}

func (m *Maze) FlagsCaptured() int {
	return m.captured
}

func (m *Maze) FlagsRemaining() int {
	return m.flags - m.captured
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

func (m *Maze) Size() (int, int) {
	return m.width, m.height
}

func (m *Maze) Robot() *Robot {
	return m.robot
}

func (m *Maze) StopRobot() {
	*m.robot = m.robot.Stop()
}
func (m *Maze) AdvanceRobot() {
	robot := m.robot.Advance()
	cell := m.CellAt(robot.X, robot.Y)
	if cell.Flag() && !cell.Captured() {
		log.Printf("capture %d %d", robot.X, robot.Y)
		m.CaptureFlag(robot.X, robot.Y)
	}
	if col := robot.ColorPainting(); col != NoColor {
		m.PaintCell(robot.X, robot.Y, col)
	}
	*m.robot = robot
}

func (m *Maze) CommandRobot(com Command) bool {
	next := m.robot.ApplyCommand(com)
	crash := next.IsMovingForward() && m.HasWallAt(next.X, next.Y, next.Orientation)
	if crash {
		next = m.robot.ApplyCommand(NoCommand)
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
