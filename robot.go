package main

import (
	"fmt"
	"math"
)

type Mover struct {
	Orientation
	Rotation
	Position
	Velocity
}

type Orientation int

const (
	North Orientation = iota
	East
	South
	West
)

func (o Orientation) String() string {
	switch o {
	case North:
		return "north"
	case East:
		return "east"
	case South:
		return "south"
	case West:
		return "west"
	default:
		return "unknown"
	}
}

type Rotation int

func (r Rotation) String() string {
	switch r {
	case Left:
		return "left"
	case Right:
		return "right"
	case NoRotation:
		return "no rotation"
	default:
		return "unknown"
	}
}

func (r Rotation) DegreesAt(t float64) float64 {
	return float64(r) * (math.Pi / 2) * t
}

const (
	Left       Rotation = -1
	Right      Rotation = 1
	NoRotation Rotation = 0
)

type Position struct {
	X, Y int
}

func (p Position) Coords() (float64, float64) {
	return float64(p.X), float64(p.Y)
}

func (p Position) String() string {
	return fmt.Sprintf("(%d, %d)", p.X, p.Y)
}

type Velocity struct {
	Dx, Dy int
}

func (v Velocity) TranslationAt(t float64) (float64, float64) {
	return float64(v.Dx) * t, float64(v.Dy) * t
}

func (v Velocity) Orientation() (Orientation, bool) {
	switch {
	case v.Dx == 0 && v.Dy > 0:
		return South, true
	case v.Dx == 0 && v.Dy < 0:
		return North, true
	case v.Dy == 0 && v.Dx > 0:
		return East, true
	case v.Dy == 0 && v.Dx < 0:
		return West, true
	default:
		return 0, false
	}
}

func (o Orientation) Rotate(r Rotation) Orientation {
	o = (o + Orientation(r)) % 4
	if o < 0 {
		return o + 4
	}
	return o
}

func (o Orientation) Reverse() Orientation {
	return (o + 2) % 4
}

func (o Orientation) VelocityForward() Velocity {
	switch o {
	case North:
		return Velocity{Dy: -1}
	case East:
		return Velocity{Dx: 1}
	case South:
		return Velocity{Dy: 1}
	case West:
		return Velocity{Dx: -1}
	default:
		return Velocity{} // or panic?
	}
}

func (o Orientation) Angle() float64 {
	return float64(o) * (math.Pi / 2)
}

func (p Position) Move(v Velocity) Position {
	return Position{
		X: p.X + v.Dx,
		Y: p.Y + v.Dy,
	}
}

type Robot struct {
	Position
	Orientation
	Velocity
	Rotation
	CurrentCommand Command
}

func (r Robot) AngleAt(t float64) float64 {
	return r.Angle() + r.Rotation.DegreesAt(t)
}

func (r Robot) CoordsAt(t float64) (float64, float64) {
	x, y := r.Coords()
	dx, dy := r.TranslationAt(t)
	return x + dx, y + dy
}

func (r Robot) ApplyCommand(com Command) Robot {
	r.Rotation = NoRotation
	r.Velocity = Velocity{}
	r.CurrentCommand = com
	switch com {
	case TurnLeft:
		r.Rotation = Left
	case TurnRight:
		r.Rotation = Right
	case MoveForward:
		r.Velocity = r.VelocityForward()
	}
	return r
}

func (r Robot) Stop() Robot {
	r.Velocity = Velocity{}
	r.Rotation = NoRotation
	return r
}

func (r Robot) Advance() Robot {
	r.Position = r.Position.Move(r.Velocity)
	r.Orientation = r.Orientation.Rotate(r.Rotation)
	return r
}

func (r Robot) IsMovingForward() bool {
	return r.Velocity == r.VelocityForward()
}

func (r Robot) ColorPainting() Color {
	switch r.CurrentCommand {
	case PaintBlue:
		return Blue
	case PaintRed:
		return Red
	case PaintYellow:
		return Yellow
	default:
		return NoColor
	}
}
