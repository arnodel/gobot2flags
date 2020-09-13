package main

import (
	"github.com/hajimehoshi/ebiten"
)

type Canvas interface {
	DrawImage(img *ebiten.Image, options *ebiten.DrawImageOptions) error
}

type transformCanvas struct {
	target   *ebiten.Image
	baseGeoM ebiten.GeoM
}

func (s *transformCanvas) DrawImage(img *ebiten.Image, options *ebiten.DrawImageOptions) error {
	op := *options
	op.GeoM.Concat(s.baseGeoM)
	return s.target.DrawImage(img, &op)
}

type MazeRenderer struct {
	cellWidth, cellHeight int
	wallWidth, wallHeight int
	walls                 Walls
	floors                Floors
	flag, robot           *Sprite
}

func (r *MazeRenderer) DrawFloor(c Canvas, x, y int, col Color) {
	op := ebiten.DrawImageOptions{}
	op.GeoM.Translate(float64(x*r.cellWidth), float64(y*r.cellHeight))
	c.DrawImage(r.floors.GetImage(col), &op)
}

func (r *MazeRenderer) DrawNorthWall(c Canvas, x, y int) {
	op := ebiten.DrawImageOptions{}
	op.GeoM.Translate(float64(x*r.cellWidth), float64(y*r.cellHeight-r.wallHeight))
	c.DrawImage(r.walls.Horizontal, &op)
}

func (r *MazeRenderer) DrawWestWall(c Canvas, x, y int) {
	op := ebiten.DrawImageOptions{}
	op.GeoM.Translate(float64(x*r.cellWidth-r.wallWidth/2), float64(y*r.cellHeight))
	c.DrawImage(r.walls.Vertical, &op)
}

func (r *MazeRenderer) DrawCornerWall(c Canvas, x, y int) {
	op := ebiten.DrawImageOptions{}
	op.GeoM.Translate(float64(x*r.cellWidth-r.wallWidth/2), float64(y*r.cellHeight-r.wallHeight))
	c.DrawImage(r.walls.Corner, &op)
}

func (r *MazeRenderer) Flag(x, y int, frame int, captured bool) ImageToDraw {
	variant := 0
	if captured {
		variant = 1
	}
	op := ebiten.DrawImageOptions{}
	tr := &op.GeoM
	r.flag.Anchor(tr)
	flagY := float64(y*r.cellHeight + 9)
	op.GeoM.Translate(float64(x*r.cellWidth+6), flagY)
	return ImageToDraw{
		Image:   r.flag.GetImage(variant, frame),
		Options: &op,
		Z:       flagY,
	}
}

func (r *MazeRenderer) Robot(robot *Robot, t float64, frame int) ImageToDraw {
	a := robot.AngleAt(t)
	x, y := robot.CoordsAt(t)
	op := ebiten.DrawImageOptions{}
	tr := &op.GeoM
	r.robot.Anchor(tr)
	tr.Rotate(a)
	robotY := (y + 0.5) * float64(r.cellHeight)
	tr.Translate((x+0.5)*float64(r.cellWidth), robotY)
	return ImageToDraw{
		Image:   r.robot.GetImage(0, 0),
		Options: &op,
		Z:       robotY,
	}
}
