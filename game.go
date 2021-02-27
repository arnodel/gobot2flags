package main

import (
	"image"
	"image/color"
	"math"

	"github.com/arnodel/gobot2flags/engine"
	"github.com/arnodel/gobot2flags/model"
	"github.com/hajimehoshi/ebiten/v2"
)

type Game struct {
	outsideWidth, outsideHeight int
	count                       int
	step                        int
	showBoard                   bool
	proportion                  float64
	mazeRenderer                *MazeRenderer
	maze                        *model.Maze
	boardRenderer               *CircuitBoardRenderer
	board                       *model.CircuitBoard
	chipSelector                *boardTiles
	boardController             *model.BoardController
	mazeWindow                  *engine.Window
	mazeControlsWindow          *engine.Window
	boardWindow                 *engine.Window
	boardControlsWindow         *engine.Window
	gameControlSelector         *gameControlSelector
	pointer                     *engine.PointerTracker
	playing                     bool
}

func (g *Game) Update() error {
	screenRect := image.Rect(0, 0, g.outsideWidth, g.outsideHeight)
	var mr, br image.Rectangle
	var tr, mtr, btr ebiten.GeoM
	const scale = 2
	if g.showBoard {
		g.proportion = math.Max(0, g.proportion-0.1)
	} else {
		g.proportion = math.Min(1, g.proportion+0.1)
	}
	if g.outsideWidth > g.outsideHeight {
		mr, br = vSplit(screenRect, int(float64(g.outsideWidth)*(1+g.proportion)/3))
	} else {
		mr, br = hSplit(screenRect, int(float64(g.outsideHeight)*(1+g.proportion)/3))
	}
	tr.Scale(2, 2)
	btr.Scale(2-g.proportion, 2-g.proportion)
	mtr.Scale(1+g.proportion, 1+g.proportion)

	// board
	br1, br2 := hSplit(br, int(128*(1-g.proportion)))

	g.boardControlsWindow = engine.CenteredWindow(br1, g.chipSelector.Bounds(), tr)
	g.boardWindow = engine.CenteredWindow(br2, g.boardRenderer.CircuitBoardBounds(g.board), btr)

	//gameWon := g.boardController.GameWon()

	// maze
	var adv int
	switch g.gameControlSelector.selectedControl {
	case NoControl:
		// Paused
	case Play, Step:
		adv = 1
	case Pause:
		// Paused
	case FastForward:
		adv = 5 - g.step%5
	case Rewind:
		if g.playing {
			g.board.ClearActiveChips()
			g.boardController = nil
			g.playing = false
			g.step = 0
		}
	}
	if !g.playing && adv > 0 {
		boardController := model.NewBoardController(g.board, g.maze.Clone())
		if boardController != nil {
			g.boardController = boardController
			g.playing = true
		}
	}
	if !g.playing {
		g.gameControlSelector.selectedControl = Rewind
	} else if g.gameControlSelector.selectedControl != Pause && g.step%60 == 0 {
		g.step = 0
		g.boardController.Advance()
	}
	if g.playing && adv > 0 {
		g.count++
		g.step += adv
		if g.step == 60 && g.gameControlSelector.selectedControl == Step {
			g.gameControlSelector.selectedControl = Pause
		}
	}

	mr1, mr2 := hSplit(mr, 64)

	g.mazeControlsWindow = engine.CenteredWindow(mr1, g.gameControlSelector.Bounds(), tr)
	g.mazeWindow = engine.CenteredWindow(mr2, g.mazeRenderer.MazeBounds(g.maze), mtr)

	g.pointer.Update()

	if g.pointer.Status() == engine.TouchDown {
		if g.boardWindow.Contains(g.pointer.CurrentPos()) && !g.showBoard {
			g.showBoard = true
			g.pointer.CancelTouch()
		} else if g.mazeWindow.Contains(g.pointer.CurrentPos()) && g.showBoard {
			g.showBoard = false
			g.pointer.CancelTouch()
		}
	}

	if !g.playing {
		g.updateBoard()
	}
	g.updateMaze()
	return nil
}

func (g *Game) updateBoard() {
	cur := g.pointer.CurrentPos()
	switch g.pointer.Status() {
	case engine.TouchDown:
		if g.boardControlsWindow.Contains(cur) {
			xx, yy := g.boardControlsWindow.Coords(cur)
			g.chipSelector.Click(xx, yy)
		} else if g.chipSelector.selectedType == model.NoChip {
			if g.boardWindow.Contains(cur) && g.chipSelector.selectedIcon == TrashCanIcon {
				g.board.Reset()
				g.chipSelector.selectedIcon = NoIcon
				g.chipSelector.selectedType = model.StartChip
			}
		} else if cx, cy, cok := g.slotCoords(cur); cok {
			newChip := g.board.ChipAt(cx, cy).WithType(g.chipSelector.selectedType)
			g.board.SetChipAt(cx, cy, newChip)
		}
	case engine.TouchUp:
		if g.chipSelector.selectedIcon != EraserIcon {
			return
		}
		cx, cy, cok := g.slotCoords(cur)
		sx, sy, sok := g.slotCoords(g.pointer.StartPos())
		if cok && sok && cx == sx && cy == sy {
			newChip := g.board.ChipAt(cx, cy).WithType(model.NoChip)
			g.board.SetChipAt(cx, cy, newChip)
		}
	case engine.Dragging:
		if g.chipSelector.selectedArrowType == model.NoArrow && g.chipSelector.selectedIcon != EraserIcon {
			return
		}
		cx, cy, cok := g.slotCoords(cur)
		lx, ly, lok := g.slotCoords(g.pointer.LastPos())
		if cok && lok {
			o, ok := model.Velocity{Dx: cx - lx, Dy: cy - ly}.Orientation()
			if ok {
				g.pointer.AdvanceStartPos() // This is so we don't erase chips by doing loops
				oldChip := g.board.ChipAt(lx, ly)
				newChip := oldChip.WithArrow(o, g.chipSelector.selectedArrowType)
				if g.chipSelector.selectedArrowType == model.ArrowNo && newChip != oldChip {
					g.chipSelector.selectedArrowType = model.ArrowYes
				}
				g.board.SetChipAt(lx, ly, newChip)

				// When erasing, also erase the other way
				if g.chipSelector.selectedIcon == EraserIcon {
					newChip := g.board.ChipAt(cx, cy).ClearArrow(o.Reverse())
					g.board.SetChipAt(cx, cy, newChip)
				}
			}
		}
	}
}

func (g *Game) slotCoords(p image.Point) (int, int, bool) {
	if !g.boardWindow.Contains(p) {
		return 0, 0, false
	}
	xx, yy := g.boardWindow.Coords(p)
	sx, sy := g.boardRenderer.GetSlotCoords(xx, yy)
	if !g.board.Contains(sx, sy) {
		return 0, 0, false
	}
	return sx, sy, true
}

func (g *Game) updateMaze() {
	switch g.pointer.Status() {
	case engine.TouchDown:
		cur := g.pointer.CurrentPos()
		if g.mazeControlsWindow.Contains(cur) {
			xx, yy := g.mazeControlsWindow.Coords(cur)
			g.gameControlSelector.Click(xx, yy)
		}
	}
}

func (g *Game) Draw(screen *ebiten.Image) {
	g.drawBoard(screen)
	g.drawMaze(screen)
	if g.playing && g.boardController.GameWon() {
		engine.DrawText(screen, "You Won!!!", 10, g.outsideHeight-10, color.RGBA{0, 255, 0, 255})
	} else {
		engine.DrawText(screen, "gobot2flags", 10, g.outsideHeight-10, color.White)
	}
}

func (g *Game) drawMaze(screen *ebiten.Image) {
	g.gameControlSelector.Draw(g.mazeControlsWindow.Canvas(screen))
	maze := g.maze
	if g.playing {
		maze = g.boardController.Maze()
	}
	g.mazeRenderer.DrawMaze(g.mazeWindow.Canvas(screen), maze, float64(g.step)/60, g.count/60)
}

func (g *Game) drawBoard(screen *ebiten.Image) {
	if !g.playing {
		g.chipSelector.Draw(g.boardControlsWindow.Canvas(screen), g.boardRenderer.chips)
	}
	g.boardRenderer.DrawCircuitBoard(g.boardWindow.Canvas(screen), g.board)
}

func hSplit(r image.Rectangle, y int) (r1 image.Rectangle, r2 image.Rectangle) {
	r1 = r
	r2 = r
	y += r.Min.Y
	r1.Max.Y = y
	r2.Min.Y = y
	return
}

func vSplit(r image.Rectangle, x int) (r1 image.Rectangle, r2 image.Rectangle) {
	r1 = r
	r2 = r
	x += r.Min.X
	r1.Max.X = x
	r2.Min.X = x
	return
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	g.outsideWidth = outsideWidth
	g.outsideHeight = outsideHeight
	return outsideWidth, outsideHeight
}
