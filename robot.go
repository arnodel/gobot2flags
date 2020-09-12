package main

import "math"

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

type Rotation int

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

type Velocity struct {
	Dx, Dy int
}

func (v Velocity) TranslationAt(t float64) (float64, float64) {
	return float64(v.Dx) * t, float64(v.Dy) * t
}

func (o Orientation) Rotate(r Rotation) Orientation {
	o = (o + Orientation(r)) % 4
	if o < 0 {
		return o + 4
	}
	return o
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

func (r Robot) Advance() Robot {
	r.Position = r.Position.Move(r.Velocity)
	r.Orientation = r.Orientation.Rotate(r.Rotation)
	return r
}

func (r Robot) IsMovingForward() bool {
	return r.Velocity == r.VelocityForward()
}
