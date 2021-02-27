package main

import (
	"image"

	"github.com/arnodel/gobot2flags/engine"
	"github.com/arnodel/gobot2flags/model"
	"github.com/arnodel/gobot2flags/sprites"
	"github.com/hajimehoshi/ebiten/v2"
)

type MazeRenderer struct {
	cellWidth, cellHeight int
	wallWidth, wallHeight int
	walls                 sprites.Walls
	floors                sprites.Floors
	flag, robot           *engine.Sprite
}

func (r *MazeRenderer) Floor(x, y int, col model.Color) engine.ImageToDraw {
	op := ebiten.DrawImageOptions{}
	op.GeoM.Translate(float64(x*r.cellWidth), float64(y*r.cellHeight))
	return engine.ImageToDraw{
		Image:   r.floors.GetImage(col),
		Options: &op,
	}
}

func (r *MazeRenderer) PaintFloor(x, y int, t float64, col model.Color) engine.ImageToDraw {
	img := r.Floor(x, y, col)
	img.Options.ColorM.Scale(1, 1, 1, t)
	return img
}

func (r *MazeRenderer) NorthWall(x, y int) engine.ImageToDraw {
	op := ebiten.DrawImageOptions{}
	op.GeoM.Translate(float64(x*r.cellWidth), float64(y*r.cellHeight-r.wallHeight))
	return engine.ImageToDraw{
		Image:   r.walls.Horizontal,
		Options: &op,
		Z:       float64(y * r.cellHeight),
	}
}

func (r *MazeRenderer) WestWall(x, y int) engine.ImageToDraw {
	op := ebiten.DrawImageOptions{}
	op.GeoM.Translate(float64(x*r.cellWidth-r.wallWidth/2), float64(y*r.cellHeight))
	return engine.ImageToDraw{
		Image:   r.walls.Vertical,
		Options: &op,
		Z:       float64((y + 1) * r.cellHeight),
	}
}

func (r *MazeRenderer) CornerWall(x, y int) engine.ImageToDraw {
	op := ebiten.DrawImageOptions{}
	op.GeoM.Translate(float64(x*r.cellWidth-r.wallWidth/2), float64(y*r.cellHeight-r.wallHeight))
	return engine.ImageToDraw{
		Image:   r.walls.Corner,
		Options: &op,
		Z:       float64(y*r.cellHeight) - 1e-3, // Subtract a small number to make it go behind the other walls
	}
}

func (r *MazeRenderer) Flag(x, y int, frame int, captured bool) engine.ImageToDraw {
	variant := 0
	if captured {
		variant = 1
	}
	op := ebiten.DrawImageOptions{}
	tr := &op.GeoM
	r.flag.Anchor(tr)
	flagY := float64(y*r.cellHeight + 9)
	op.GeoM.Translate(float64(x*r.cellWidth+6), flagY)
	return engine.ImageToDraw{
		Image:   r.flag.GetImage(variant, frame),
		Options: &op,
		Z:       flagY,
	}
}

func (r *MazeRenderer) Robot(robot *model.Robot, t float64, frame int) engine.ImageToDraw {
	a := robot.AngleAt(t)
	x, y := robot.CoordsAt(t)
	op := ebiten.DrawImageOptions{}
	tr := &op.GeoM
	r.robot.Anchor(tr)
	tr.Rotate(a)
	robotY := (y + 0.5) * float64(r.cellHeight)
	tr.Translate((x+0.5)*float64(r.cellWidth), robotY)
	return engine.ImageToDraw{
		Image:   r.robot.GetImage(0, 0),
		Options: &op,
		Z:       robotY,
	}
}

func (r *MazeRenderer) MazeBounds(m *model.Maze) image.Rectangle {
	w, h := m.Size()
	return image.Rect(-16, -16, 32*w+16, 32*h+16)
}

func (r *MazeRenderer) DrawMaze(c engine.Canvas, m *model.Maze, t float64, frame int) {
	w, h := m.Size()

	// Draw the floors first as they are under everything
	for y := 0; y < h; y++ {
		for x := 0; x < w; x++ {
			c.Draw(r.Floor(x, y, m.CellAt(x, y).Color()))
		}
	}

	stack := engine.ImageStack{}

	// Draw the walls and flags
	for y := 0; y <= h; y++ {
		for x := 0; x <= w; x++ {
			cell := m.CellAt(x, y)
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
	}

	// Draw the robot
	robot := m.Robot()
	if robot != nil {
		stack.Add(r.Robot(robot, t, frame))
		if col := robot.ColorPainting(); col != model.NoColor {
			stack.Add(r.PaintFloor(robot.X, robot.Y, t, col))
		}
	}

	stack.Draw(c)
	stack.Empty() // Reuse the underlying slice, same number of objects each time!
}
