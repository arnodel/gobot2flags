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

func (r *MazeRenderer) DrawFlag(c Canvas, x, y int, frame int, captured bool) {
	variant := 0
	if captured {
		variant = 1
	}
	op := ebiten.DrawImageOptions{}
	op.GeoM.Translate(float64(x*r.cellWidth), float64(y*r.cellHeight))
	c.DrawImage(r.flag.GetImage(variant, frame), &op)
}

func (r *MazeRenderer) DrawRobot(c Canvas, robot *Robot, t float64, frame int) {
	a := robot.AngleAt(t)
	x, y := robot.CoordsAt(t)
	tr := r.robot.Rotate(a)
	tr.Translate(x*float64(r.cellWidth), y*float64(r.cellHeight))
	op := ebiten.DrawImageOptions{}
	op.GeoM = tr
	c.DrawImage(r.robot.GetImage(0, 0), &op)
}
