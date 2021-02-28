package play

import (
	"image"
	"image/color"
	"math"

	"github.com/arnodel/gobot2flags/engine"
	"github.com/arnodel/gobot2flags/model"
	"github.com/arnodel/gobot2flags/sprites"
	"github.com/hajimehoshi/ebiten/v2"
)

type View struct {
	count               int
	step                int
	showBoard           bool
	proportion          float64
	mazeRenderer        *MazeRenderer
	maze                *model.Maze
	boardRenderer       *CircuitBoardRenderer
	board               *model.CircuitBoard
	chipSelector        *boardTiles
	boardController     *model.LevelController
	mazeWindow          *engine.Window
	mazeControlsWindow  *engine.Window
	boardWindow         *engine.Window
	boardControlsWindow *engine.Window
	exitWindow          *engine.Window
	gameControlSelector *gameControlSelector
	playing             bool
	exit                func()
}

var _ engine.View = (*View)(nil)

func NewView(level *model.Level, exit func()) *View {
	maze := level.Maze
	board := model.NewCircuitBoard(level.BoardWidth, level.BoardHeigth)
	chips := ChipRenderer{sprites.CircuitBoardTiles}
	boardRenderer := NewCircuitBoardRenderer(chips)
	mazeRenderer := &MazeRenderer{
		cellWidth:  sprites.FrameWidth,
		cellHeight: sprites.FrameHeight,
		wallWidth:  sprites.WallWidth,
		wallHeight: sprites.WallHeight,
		walls:      sprites.GreyWalls,
		floors:     sprites.PlainFloors,
		robot:      sprites.Robot,
		flag:       sprites.Flag,
	}
	return &View{
		maze:            maze,
		mazeRenderer:    mazeRenderer,
		board:           board,
		boardRenderer:   &boardRenderer,
		boardController: model.NewLevelController(board, maze.Clone()),
		showBoard:       true,
		chipSelector: &boardTiles{
			selectedType: model.StartChip,
			icons:        sprites.PlainIcons,
		},
		gameControlSelector: &gameControlSelector{
			selectedControl: Rewind,
			icons:           sprites.PlainIcons,
		},
		exit: exit,
	}
}

func (v *View) Update(vc engine.ViewContainer) error {
	outsideWidth, outsideHeight := vc.OutsideSize()
	screenRect := vc.OutsideRect()
	var mr, br image.Rectangle
	var tr, mtr, btr ebiten.GeoM
	const scale = 2
	if v.showBoard {
		v.proportion = math.Max(0, v.proportion-0.1)
	} else {
		v.proportion = math.Min(1, v.proportion+0.1)
	}
	if outsideWidth > outsideHeight {
		mr, br = vSplit(screenRect, int(float64(outsideWidth)*(1+v.proportion)/3))
	} else {
		mr, br = hSplit(screenRect, int(float64(outsideHeight)*(1+v.proportion)/3))
	}
	tr.Scale(2, 2)
	btr.Scale(2-v.proportion, 2-v.proportion)
	mtr.Scale(1+v.proportion, 1+v.proportion)

	// board
	br1, br2 := hSplit(br, int(128*(1-v.proportion)))

	v.boardControlsWindow = engine.CenteredWindow(br1, v.chipSelector.Bounds(), tr)
	v.boardWindow = engine.CenteredWindow(br2, v.boardRenderer.CircuitBoardBounds(v.board), btr)

	// maze
	var adv int
	switch v.gameControlSelector.selectedControl {
	case NoControl:
		// Paused
	case Play, Step:
		adv = 1
	case Pause:
		// Paused
	case FastForward:
		adv = 5 - v.step%5
	case Rewind:
		if v.playing {
			v.board.ClearActiveChips()
			v.boardController = nil
			v.playing = false
			v.step = 0
		}
	}
	if !v.playing && adv > 0 {
		boardController := model.NewLevelController(v.board, v.maze.Clone())
		if boardController != nil {
			v.boardController = boardController
			v.playing = true
		}
	}
	if !v.playing {
		v.gameControlSelector.selectedControl = Rewind
	} else if v.gameControlSelector.selectedControl != Pause && v.step%60 == 0 {
		v.step = 0
		v.boardController.Advance()
	}
	if v.playing && adv > 0 {
		v.count++
		v.step += adv
		if v.step == 60 && v.gameControlSelector.selectedControl == Step {
			v.gameControlSelector.selectedControl = Pause
		}
	}

	mr1, mr2 := hSplit(mr, 64)
	mr11, mr12 := vSplit(mr1, 32)

	v.exitWindow = engine.CenteredWindow(mr11, sprites.PlainIcons.Bounds(), ebiten.GeoM{})
	v.mazeControlsWindow = engine.CenteredWindow(mr12, v.gameControlSelector.Bounds(), tr)
	v.mazeWindow = engine.CenteredWindow(mr2, v.mazeRenderer.MazeBounds(v.maze), mtr)

	pointer := vc.Pointer()

	if pointer.Status() == engine.TouchDown {
		if v.exitWindow.Contains(pointer.CurrentPos()) {
			v.exit()
		}
		if v.boardWindow.Contains(pointer.CurrentPos()) && !v.showBoard {
			v.showBoard = true
			pointer.CancelTouch()
		} else if v.mazeWindow.Contains(pointer.CurrentPos()) && v.showBoard {
			v.showBoard = false
			pointer.CancelTouch()
		}
	}

	if !v.playing {
		v.updateBoard(pointer)
	}
	v.updateMaze(pointer)
	return nil
}

func (g *View) updateBoard(pointer *engine.PointerTracker) {
	cur := pointer.CurrentPos()
	switch pointer.Status() {
	case engine.TouchDown:
		if g.boardControlsWindow.Contains(cur) {
			xx, yy := g.boardControlsWindow.Coords(cur)
			g.chipSelector.Click(xx, yy)
		} else if g.chipSelector.selectedType == model.NoChip {
			if g.boardWindow.Contains(cur) && g.chipSelector.selectedIcon == sprites.TrashCanIcon {
				g.board.Reset()
				g.chipSelector.selectedIcon = sprites.NoIcon
				g.chipSelector.selectedType = model.StartChip
			}
		} else if cx, cy, cok := g.slotCoords(cur); cok {
			newChip := g.board.ChipAt(cx, cy).WithType(g.chipSelector.selectedType)
			g.board.SetChipAt(cx, cy, newChip)
		}
	case engine.TouchUp:
		if g.chipSelector.selectedIcon != sprites.EraserIcon {
			return
		}
		cx, cy, cok := g.slotCoords(cur)
		sx, sy, sok := g.slotCoords(pointer.StartPos())
		if cok && sok && cx == sx && cy == sy {
			newChip := g.board.ChipAt(cx, cy).WithType(model.NoChip)
			g.board.SetChipAt(cx, cy, newChip)
		}
	case engine.Dragging:
		if g.chipSelector.selectedArrowType == model.NoArrow && g.chipSelector.selectedIcon != sprites.EraserIcon {
			return
		}
		cx, cy, cok := g.slotCoords(cur)
		lx, ly, lok := g.slotCoords(pointer.LastPos())
		if cok && lok {
			o, ok := model.Velocity{Dx: cx - lx, Dy: cy - ly}.Orientation()
			if ok {
				pointer.AdvanceStartPos() // This is so we don't erase chips by doing loops
				oldChip := g.board.ChipAt(lx, ly)
				newChip := oldChip.WithArrow(o, g.chipSelector.selectedArrowType)
				if g.chipSelector.selectedArrowType == model.ArrowNo && newChip != oldChip {
					g.chipSelector.selectedArrowType = model.ArrowYes
				}
				g.board.SetChipAt(lx, ly, newChip)

				// When erasing, also erase the other way
				if g.chipSelector.selectedIcon == sprites.EraserIcon {
					newChip := g.board.ChipAt(cx, cy).ClearArrow(o.Reverse())
					g.board.SetChipAt(cx, cy, newChip)
				}
			}
		}
	}
}

func (g *View) slotCoords(p image.Point) (int, int, bool) {
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

func (g *View) updateMaze(pointer *engine.PointerTracker) {
	switch pointer.Status() {
	case engine.TouchDown:
		cur := pointer.CurrentPos()
		if g.mazeControlsWindow.Contains(cur) {
			xx, yy := g.mazeControlsWindow.Coords(cur)
			g.gameControlSelector.Click(xx, yy)
		}
	}
}

func (g *View) Draw(screen *ebiten.Image) {
	g.exitWindow.Canvas(screen).Draw(sprites.PlainIcons.ImageToDraw(sprites.BackIcon))
	g.drawBoard(screen)
	g.drawMaze(screen)
	maxY := screen.Bounds().Max.Y
	if g.playing && g.boardController.GameWon() {
		engine.DrawText(screen, "You Won!!!", 10, maxY-10, color.RGBA{0, 255, 0, 255})
	}
}

func (g *View) drawMaze(screen *ebiten.Image) {
	g.gameControlSelector.Draw(g.mazeControlsWindow.Canvas(screen))
	maze := g.maze
	if g.playing {
		maze = g.boardController.Maze()
	}
	g.mazeRenderer.DrawMaze(g.mazeWindow.Canvas(screen), maze, float64(g.step)/60, g.count/60)
}

func (g *View) drawBoard(screen *ebiten.Image) {
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
