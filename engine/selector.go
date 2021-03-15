package engine

import "math"

type SelectorStatus int

const (
	NotSelecting SelectorStatus = iota
	Selecting
	Select
)

type Selector struct {
	Status         SelectorStatus
	SelectingIndex int
	SelectIndex    int
	SelectingTicks int
}

func (s *Selector) Update(i int, t TouchStatus) SelectorStatus {
	s.SelectIndex = 0
	if i < 0 {
		t = NoTouch
	}
	switch t {
	case TouchDown:
		s.SelectingIndex = i
		s.Status = Selecting
		s.SelectingTicks = 0
	case TouchUp:
		if s.Status == Selecting && s.SelectingIndex == i {
			s.SelectIndex = i
			s.Status = Select
		} else {
			s.SelectingIndex = 0
			s.Status = NotSelecting
		}
		s.SelectingTicks = 0
	case Dragging:
		if s.SelectingIndex != i {
			s.Status = NotSelecting
			s.SelectingIndex = 0
			s.SelectingTicks = 0
		} else {
			s.SelectingTicks++
		}
	default:
		s.SelectingIndex = 0
		s.Status = NotSelecting
		s.SelectingTicks = 0
	}
	return s.Status
}

func (s *Selector) IsSelecting(i int) bool {
	return s.Status == Selecting && s.SelectingIndex == i
}

func (s *Selector) SelectingProportion() float64 {
	return math.Min(float64(s.SelectingTicks)/10, 1)
}
