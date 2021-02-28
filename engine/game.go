package engine

import (
	"image"
	"sync"

	"github.com/hajimehoshi/ebiten/v2"
)

type View interface {
	Update(ViewContainer) error
	Draw(*ebiten.Image)
}

type ViewContainer interface {
	OutsideSize() (int, int)
	OutsideRect() image.Rectangle
	Pointer() *PointerTracker
}

type Game struct {
	outsideWidth  int
	outsideHeight int
	currentView   View
	nextView      View
	mx            sync.Mutex
	pointer       PointerTracker
}

var _ ebiten.Game = (*Game)(nil)
var _ ViewContainer = (*Game)(nil)

func NewGame(initialView View) *Game {
	return &Game{
		currentView: initialView,
	}
}

func (g *Game) Update() error {
	g.pointer.Update()
	g.mx.Lock()
	if g.nextView != nil {
		g.currentView = g.nextView
		g.nextView = nil
		g.pointer.CancelTouch()
	}
	g.mx.Unlock()
	if g.currentView != nil {
		return g.currentView.Update(g)
	}
	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	if g.currentView != nil {
		g.currentView.Draw(screen)
	}
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	g.outsideWidth = outsideWidth
	g.outsideHeight = outsideHeight
	return outsideWidth, outsideHeight
}

func (g *Game) OutsideSize() (int, int) {
	return g.outsideWidth, g.outsideHeight
}

func (g *Game) OutsideRect() image.Rectangle {
	return image.Rect(0, 0, g.outsideWidth, g.outsideHeight)
}

func (g *Game) Pointer() *PointerTracker {
	return &g.pointer
}

func (g *Game) SetView(v View) {
	g.mx.Lock()
	g.nextView = v
	g.mx.Unlock()
}
