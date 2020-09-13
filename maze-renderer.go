package main

import (
	"sort"

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

type imageToDraw struct {
	image   *ebiten.Image
	options *ebiten.DrawImageOptions
	y       float64
}

type MazeRenderer struct {
	cellWidth, cellHeight int
	wallWidth, wallHeight int
	walls                 Walls
	floors                Floors
	flag, robot           *Sprite
	pending               []imageToDraw
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

func (r *MazeRenderer) addPending(img *ebiten.Image, op *ebiten.DrawImageOptions, y float64) {
	r.pending = append(r.pending, imageToDraw{
		image:   img,
		options: op,
		y:       y,
	})
}

func (r *MazeRenderer) DrawPending(c Canvas) {
	imgs := r.pending
	sort.Slice(imgs, func(i, j int) bool {
		return imgs[i].y < imgs[j].y
	})
	for _, toDraw := range imgs {
		c.DrawImage(toDraw.image, toDraw.options)
	}
	r.pending = nil
}

func (r *MazeRenderer) AddFlag(x, y int, frame int, captured bool) {
	variant := 0
	if captured {
		variant = 1
	}
	op := ebiten.DrawImageOptions{}
	tr := &op.GeoM
	r.flag.Anchor(tr)
	flagY := float64(y*r.cellHeight + 9)
	op.GeoM.Translate(float64(x*r.cellWidth+6), flagY)
	r.addPending(r.flag.GetImage(variant, frame), &op, flagY)
}

func (r *MazeRenderer) AddRobot(robot *Robot, t float64, frame int) {
	a := robot.AngleAt(t)
	x, y := robot.CoordsAt(t)
	op := ebiten.DrawImageOptions{}
	tr := &op.GeoM
	r.robot.Anchor(tr)
	tr.Rotate(a)
	robotY := (y + 0.5) * float64(r.cellHeight)
	tr.Translate((x+0.5)*float64(r.cellWidth), robotY)
	r.addPending(r.robot.GetImage(0, 0), &op, robotY)
}
