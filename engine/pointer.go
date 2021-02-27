package engine

import (
	"image"

	"github.com/hajimehoshi/ebiten/v2"
)

type PointerStatus int

const (
	NoTouch PointerStatus = iota
	TouchDown
	Dragging
	TouchUp
)

type PointerTracker struct {
	startPos    image.Point
	lastPos     image.Point
	currentPos  image.Point
	status      PointerStatus
	frames      int
	cancelTouch bool
}

func (p *PointerTracker) CancelTouch() {
	p.cancelTouch = true
	p.status = NoTouch
}

func (p *PointerTracker) Status() PointerStatus {
	return p.status
}

func (p *PointerTracker) CurrentPos() image.Point {
	return p.currentPos
}

func (p *PointerTracker) LastPos() image.Point {
	return p.lastPos
}

func (p *PointerTracker) StartPos() image.Point {
	return p.startPos
}

// TODO: perhaps remove necessity for this method?
func (p *PointerTracker) AdvanceStartPos() {
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
