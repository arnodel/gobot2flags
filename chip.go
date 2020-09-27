package main

type Chip uint16

func (c Chip) Type() ChipType {
	return ChipType(c & 0xff)
}

func (c Chip) IsTest() bool {
	return c&0x10 != 0
}
func (c Chip) ArrowYes() (Orientation, bool) {
	data := (c >> 8) & 0xf
	return Orientation(data & 0x3), (data & 0x4) != 0
}

func (c Chip) ArrowNo() (Orientation, bool) {
	data := c >> 12
	return Orientation(data & 0x3), (data & 0x4) != 0
}

func (c Chip) Arrow(ok bool) (Orientation, bool) {
	if ok {
		return c.ArrowYes()
	}
	return c.ArrowNo()
}

func (c Chip) WithType(t ChipType) Chip {
	return (c &^ 0xff) | Chip(t)
}

func (c Chip) ClearArrowYes() Chip {
	return c &^ 0xf00
}

func (c Chip) ClearArrowNo() Chip {
	return c &^ 0xf000
}

func (c Chip) WithArrowYes(o Orientation) Chip {
	return c.ClearArrowYes() | Chip((0x4|o)<<8)
}

func (c Chip) WithArrowNo(o Orientation) Chip {
	return c.ClearArrowNo() | Chip((0x4|o)<<12)
}

func (c Chip) Command(floorColor Color, wallAhead bool) (Command, bool) {
	switch c.Type() {
	case StartChip:
		return NoCommand, true
	case ForwardChip:
		return MoveForward, true
	case TurnLeftChip:
		return TurnLeft, true
	case TurnRightChip:
		return TurnRight, true
	case PaintRedChip:
		return PaintRed, true
	case PaintBlueChip:
		return PaintBlue, true
	case PaintYellowChip:
		return PaintYellow, true
	case IsWallAheadChip:
		return NoCommand, wallAhead
	case IsFloorRedChip:
		return NoCommand, floorColor == Red
	case IsFloorBlueChip:
		return NoCommand, floorColor == Blue
	case IsFloorYellowChip:
		return NoCommand, floorColor == Yellow
	default:
		return NoCommand, false
	}
}

type ChipType byte

const (
	NoChip ChipType = iota
	StartChip
	ForwardChip
	TurnLeftChip
	TurnRightChip
	PaintRedChip
	PaintYellowChip
	PaintBlueChip

	IsWallAheadChip = iota + 0x10
	IsFloorRedChip
	IsFloorYellowChip
	IsFloorBlueChip
)
