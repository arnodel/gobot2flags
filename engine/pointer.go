package engine

import (
	"image"

	"github.com/hajimehoshi/ebiten/v2"
)

type TouchStatus int

const (
	NoTouch TouchStatus = iota
	TouchDown
	Dragging
	TouchUp
)

type PointerTracker struct {
	startPos    image.Point
	lastPos     image.Point
	currentPos  image.Point
	status      TouchStatus
	frames      int
	cancelTouch bool
}

func (p *PointerTracker) CancelTouch() {
	if p == nil {
		return
	}
	p.cancelTouch = true
	p.status = NoTouch
}

func (p *PointerTracker) Status() TouchStatus {
	if p == nil {
		return NoTouch
	}
	return p.status
}

func (p *PointerTracker) CurrentPos() image.Point {
	if p == nil {
		return image.Point{}
	}
	return p.currentPos
}

func (p *PointerTracker) LastPos() image.Point {
	if p == nil {
		return image.Point{}
	}
	return p.lastPos
}

func (p *PointerTracker) StartPos() image.Point {
	if p == nil {
		return image.Point{}
	}
	return p.startPos
}

// TODO: perhaps remove necessity for this method?
func (p *PointerTracker) AdvanceStartPos() {
	if p == nil {
		return
	}
	p.startPos = p.lastPos
}

func (p *PointerTracker) Update() {
	var currentPos image.Point
	if ebiten.IsMouseButtonPressed(ebiten.MouseButtonLeft) {
		currentPos = image.Pt(ebiten.CursorPosition())
	} else if touchIDs := ebiten.TouchIDs(); touchIDs != nil {
		currentPos = image.Pt(ebiten.TouchPosition(touchIDs[0]))
	} else {
		p.cancelTouch = false
		switch p.status {
		case NoTouch:
			// Nothing to do?
		case TouchDown, Dragging:
			p.status = TouchUp
			p.lastPos = p.currentPos
			p.frames = 0
		case TouchUp:
			p.status = NoTouch
		}
		return
	}
	if p.cancelTouch {
		return
	}
	switch p.status {
	case NoTouch, TouchUp:
		p.status = TouchDown
		p.startPos = currentPos
		p.lastPos = currentPos
		p.frames = 1
	case TouchDown:
		p.status = Dragging
		fallthrough
	case Dragging:
		if currentPos == p.lastPos && p.frames > 0 {
			p.frames++
		} else {
			p.lastPos = p.currentPos
			p.frames = 0
		}
	default:
		// Shouldn't get there?
	}
	p.currentPos = currentPos
	return
}

func (p *PointerTracker) ForWindow(w *Window) PointerStatus {
	if !w.Contains(p.currentPos) {
		return PointerStatus{}
	}
	return PointerStatus{
		tracker: p,
		tr:      w.inv,
	}
}

type PointerStatus struct {
	tracker *PointerTracker
	tr      ebiten.GeoM
}

func (s PointerStatus) HasPointer() bool {
	return s.tracker != nil
}

func (s PointerStatus) Status() TouchStatus {
	return s.tracker.Status()
}

func (s PointerStatus) CurrentCoords() (float64, float64) {
	return s.coords(s.tracker.CurrentPos())
}

func (s PointerStatus) LastCoords() (float64, float64) {
	return s.coords(s.tracker.LastPos())
}

func (s PointerStatus) StartCoords() (float64, float64) {
	return s.coords(s.tracker.StartPos())
}

func (s PointerStatus) CancelTouch() {
	s.tracker.CancelTouch()
}

func (s PointerStatus) AdvanceStartPos() {
	s.tracker.AdvanceStartPos()
}

func (s PointerStatus) coords(pt image.Point) (float64, float64) {
	return s.tr.Apply(fc(pt))
}
