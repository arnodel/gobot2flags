package main

import (
	"image"
	_ "image/png"
	"log"

	"github.com/arnodel/gobots2flags/resources"
	"github.com/hajimehoshi/ebiten/v2"
)

const (
	screenWidth  = 320
	screenHeight = 240

	frameOX     = 0
	frameOY     = 0
	frameWidth  = 32
	frameHeight = 32
	frameNum    = 6

	wallWidth  = 6
	wallHeight = 7
)

type PointerStatus int

const (
	NoTouch PointerStatus = iota
	TouchDown
	Dragging
	TouchUp
)

type PointerTracker struct {
	startPos   image.Point
	lastPos    image.Point
	currentPos image.Point
	status     PointerStatus
	frames     int
}

func (p *PointerTracker) Update() {
	x, y := ebiten.CursorPosition()
	currentPos := image.Point{x, y}
	if ebiten.IsMouseButtonPressed(ebiten.MouseButtonLeft) {
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
	} else {
		switch p.status {
		case NoTouch:
			// Nothing to do?
		case TouchDown, Dragging:
			p.status = TouchUp
			if currentPos == p.lastPos && p.frames > 0 {
				p.frames++
			} else {
				p.lastPos = p.currentPos
				p.frames = 0
			}
		case TouchUp:
			p.status = NoTouch
		}
	}
}

type Game struct {
	outsideWidth, outsideHeight int
	count                       int
	step                        int
	showBoard                   bool
	mazeRenderer                *MazeRenderer
	maze                        *Maze
	boardRenderer               *CircuitBoardRenderer
	board                       *CircuitBoard
	chipSelector                *boardTiles
	controller                  *ManualController
	boardController             *BoardController
	mainWindow                  *Window
	controlsWindow              *Window
	gameControlSelector         *gameControlSelector
	pointer                     PointerTracker
}

func (g *Game) Update() error {
	g.controller.UpdateNextCommand()
	if g.controller.GetSwitch() {
		g.showBoard = !g.showBoard
	}
	screenRect := image.Rect(0, 0, g.outsideWidth, g.outsideHeight)
	if g.showBoard {
		r1, r2 := hSplit(screenRect, 64)

		const scale = 2

		baseTr1 := ebiten.GeoM{}
		baseTr1.Translate(-float64((len(arrowTypes)+len(chipTypes))*24)/2, -float64(frameHeight)/2)
		baseTr1.Scale(scale, scale)

		baseTr2 := ebiten.GeoM{}
		baseTr2.Translate(-float64(g.board.width*frameWidth)/2, -float64(g.board.height*frameHeight)/2)
		baseTr2.Scale(scale, scale)

		g.controlsWindow = &Window{
			bounds: r1,
			tr:     baseTr1,
		}
		g.controlsWindow.Center()
		g.mainWindow = &Window{
			bounds: r2,
			tr:     baseTr2,
		}
		g.mainWindow.Center()
	} else {
		var adv int
		switch g.gameControlSelector.selectedControl {
		case NoControl:
			// Paused
		case Play:
			adv = 1
		case FastForward:
			adv = 5 - g.step%5
		case Rewind:
			g.boardController = newBoardController(g.board, g.maze.Clone())
			g.gameControlSelector.selectedControl = NoControl
		}
		if adv > 0 {
			g.count++
			g.step = (g.step + adv) % 60
			if g.step == 0 {
				g.boardController.maze.AdvanceRobot(g.boardController.NextCommand())
			}
		}

		r1, r2 := hSplit(screenRect, 64)
		const scale = 2

		baseTr1 := ebiten.GeoM{}
		baseTr1.Translate(-float64(len(gameControls)*24)/2, -float64(frameHeight)/2)
		baseTr1.Scale(scale, scale)

		baseTr2 := ebiten.GeoM{}
		baseTr2.Translate(-float64(g.maze.width*frameWidth)/2, -float64(g.maze.height*frameHeight)/2)
		baseTr2.Scale(scale, scale)

		g.controlsWindow = &Window{
			bounds: r1,
			tr:     baseTr1,
		}
		g.controlsWindow.Center()
		g.mainWindow = &Window{
			bounds: r2,
			tr:     baseTr2,
		}
		g.mainWindow.Center()
	}
	g.pointer.Update()
	if g.showBoard {
		switch g.pointer.status {
		case TouchUp:
			if g.pointer.frames > 0 {
				x, y := g.pointer.currentPos.X, g.pointer.currentPos.Y
				if g.controlsWindow.Contains(x, y) {
					xx, yy := g.controlsWindow.Coords(x, y)
					g.chipSelector.Click(xx, yy)
				} else if g.mainWindow.Contains(x, y) {
					xx, yy := g.mainWindow.Coords(x, y)
					sx, sy := g.boardRenderer.GetSlotCoords(xx, yy)
					if g.board.Contains(sx, sy) {
						newChip := g.board.ChipAt(sx, sy).WithType(g.chipSelector.selectedType)
						g.board.SetChipAt(sx, sy, newChip)
					}
				}
			}
		case Dragging:
			x1, y1 := g.pointer.lastPos.X, g.pointer.lastPos.Y
			x2, y2 := g.pointer.currentPos.X, g.pointer.currentPos.Y
			if g.mainWindow.Contains(x1, y1) && g.mainWindow.Contains(x2, y2) {
				xx1, yy1 := g.mainWindow.Coords(x1, y1)
				xx2, yy2 := g.mainWindow.Coords(x2, y2)
				sx1, sy1 := g.boardRenderer.GetSlotCoords(xx1, yy1)
				sx2, sy2 := g.boardRenderer.GetSlotCoords(xx2, yy2)
				if g.board.Contains(sx1, sy1) && g.board.Contains(sx2, sy2) {
					o, ok := Velocity{Dx: sx2 - sx1, Dy: sy2 - sy1}.Orientation()
					if ok {
						oldChip := g.board.ChipAt(sx1, sy1)
						newChip := oldChip.WithArrow(o, g.chipSelector.selectedArrowType)
						if g.chipSelector.selectedArrowType == ArrowNo && newChip != oldChip {
							g.chipSelector.selectedArrowType = ArrowYes
						}
						g.board.SetChipAt(sx1, sy1, newChip)
					}
				}
			}
		}
	} else {
		switch g.pointer.status {
		case TouchUp:
			if g.pointer.frames > 0 {
				x, y := g.pointer.currentPos.X, g.pointer.currentPos.Y
				if g.controlsWindow.Contains(x, y) {
					xx, yy := g.controlsWindow.Coords(x, y)
					g.gameControlSelector.Click(xx, yy)
				}
			}
		}
	}
	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	if g.showBoard {
		g.drawBoard(screen)
	} else {
		g.drawMaze(screen)
	}
}

func (g *Game) drawMaze(screen *ebiten.Image) {
	g.gameControlSelector.Draw(g.controlsWindow.Canvas(screen))
	g.boardController.maze.Draw(g.mainWindow.Canvas(screen), g.mazeRenderer, float64(g.step)/60, g.count/60)
}

func (g *Game) drawBoard(screen *ebiten.Image) {
	g.chipSelector.Draw(g.controlsWindow.Canvas(screen), g.boardRenderer.chips)
	g.board.Draw(g.mainWindow.Canvas(screen), g.boardRenderer)
}

func hSplit(r image.Rectangle, y int) (r1 image.Rectangle, r2 image.Rectangle) {
	r1 = r
	r2 = r
	r1.Max.Y = y
	r2.Min.Y = y
	return
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	g.outsideWidth = outsideWidth
	g.outsideHeight = outsideHeight
	return outsideWidth, outsideHeight
}

func main() {
	maze, err := MazeFromString(`
+--+--+--+--+
|RF|R |R  RF|
+  .  .  .  +
|Y  B> Y  B |
+--+--+  +  +
|BF Y  B |YF|
+--+--+--+--+
`)
	if err != nil {
		log.Fatalf("Could not create maze: %s", err)
	}
	board, err := CircuitBoardFromString(`
|ST -> W? y> TL|
| ^    nv     v|
|..    MF <- ..|
| ^     v      |
|.. <- PR      |`)
	if err != nil {
		log.Fatalf("Could not create circuit board: %s", err)
	}
	chips := ChipRenderer{NewSprite(resources.GetImage("circuitboardtiles.png"), 32, 32, 16, 16)}
	boardRenderer := NewCircuitBoardRenderer(chips)
	mazeRenderer := &MazeRenderer{
		cellWidth:  frameWidth,
		cellHeight: frameHeight,
		wallWidth:  wallWidth,
		wallHeight: wallHeight,
		walls:      NewWalls(resources.GetImage("greywalls.png")),
		floors:     LoadFloors(resources.GetImage("floors.png")),
		robot:      NewSprite(resources.GetImage("robot.png"), frameWidth, frameHeight, 16, 16),
		flag:       NewSprite(resources.GetImage("greenflag.png"), frameWidth, frameHeight, 10, 28),
	}
	game := &Game{
		maze:            maze,
		mazeRenderer:    mazeRenderer,
		board:           board,
		boardRenderer:   &boardRenderer,
		controller:      &ManualController{},
		boardController: newBoardController(board, maze.Clone()),
		showBoard:       true,
		chipSelector:    new(boardTiles),
		gameControlSelector: &gameControlSelector{
			controlImages: NewSprite(resources.GetImage("gamecontrols.png"), 32, 32, 16, 16),
		},
	}
	ebiten.SetWindowSize(screenWidth*2, screenHeight*2)
	ebiten.SetWindowTitle("Gobot 2 Flags")
	ebiten.SetWindowResizable(true)
	if err := ebiten.RunGame(game); err != nil {
		log.Fatal(err)
	}
}
