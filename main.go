package main

import (
	"flag"
	"image"
	_ "image/png"
	"io/ioutil"
	"log"
	"math"

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
	proportion                  float64
	mazeRenderer                *MazeRenderer
	maze                        *Maze
	boardRenderer               *CircuitBoardRenderer
	board                       *CircuitBoard
	chipSelector                *boardTiles
	controller                  *ManualController
	boardController             *BoardController
	mazeWindow                  *Window
	mazeControlsWindow          *Window
	boardWindow                 *Window
	boardControlsWindow         *Window
	gameControlSelector         *gameControlSelector
	pointer                     PointerTracker
	playing                     bool
}

func (g *Game) Update() error {
	g.controller.UpdateNextCommand()
	if g.controller.GetSwitch() {
		g.showBoard = !g.showBoard
	}
	screenRect := image.Rect(0, 0, g.outsideWidth, g.outsideHeight)
	var mr, br image.Rectangle
	var mtr, btr ebiten.GeoM
	const scale = 2
	if g.showBoard {
		g.proportion = math.Max(0, g.proportion-0.1)
	} else {
		g.proportion = math.Min(1, g.proportion+0.1)
	}
	mr, br = vSplit(screenRect, int(float64(g.outsideWidth)*(1+g.proportion)/3))
	btr.Scale(2-g.proportion, 2-g.proportion)
	mtr.Scale(1+g.proportion, 1+g.proportion)

	// board
	br1, br2 := hSplit(br, 64)

	g.boardControlsWindow = centeredWindow(br1, g.chipSelector.Bounds(), btr)
	g.boardWindow = centeredWindow(br2, g.board.Bounds(), btr)

	// maze
	var adv int
	switch g.gameControlSelector.selectedControl {
	case NoControl:
		// Paused
	case Play:
		adv = 1
		g.playing = true
	case FastForward:
		adv = 5 - g.step%5
		g.playing = true
	case Rewind:
		g.board.ClearActiveChips()
		g.boardController = newBoardController(g.board, g.maze.Clone())
		g.gameControlSelector.selectedControl = NoControl
		g.playing = false
	}
	if adv > 0 {
		g.count++
		g.step = (g.step + adv) % 60
		if g.step == 0 {
			g.boardController.maze.AdvanceRobot(g.boardController.NextCommand())
		}
	}

	mr1, mr2 := hSplit(mr, 64)

	g.mazeControlsWindow = centeredWindow(mr1, g.gameControlSelector.Bounds(), mtr)
	g.mazeWindow = centeredWindow(mr2, g.maze.Bounds(), mtr)

	g.pointer.Update()

	consumedTouchDown := false
	if g.pointer.status == TouchDown {
		if g.boardWindow.Contains(g.pointer.currentPos) && !g.showBoard {
			g.showBoard = true
			consumedTouchDown = true
		} else if g.mazeWindow.Contains(g.pointer.currentPos) && g.showBoard {
			g.showBoard = false
			consumedTouchDown = true
		}
	}

	if !consumedTouchDown && adv == 0 {
		g.updateBoard()
	}
	if !consumedTouchDown {
		g.updateMaze()
	}
	return nil
}

func (g *Game) updateBoard() {
	switch g.pointer.status {
	case TouchUp:

		if g.pointer.frames > 0 {
			cur := g.pointer.currentPos
			if g.boardControlsWindow.Contains(cur) {
				xx, yy := g.boardControlsWindow.Coords(cur)
				g.chipSelector.Click(xx, yy)
			} else if g.boardWindow.Contains(cur) {
				xx, yy := g.boardWindow.Coords(cur)
				sx, sy := g.boardRenderer.GetSlotCoords(xx, yy)
				if g.board.Contains(sx, sy) {
					newChip := g.board.ChipAt(sx, sy).WithType(g.chipSelector.selectedType)
					g.board.SetChipAt(sx, sy, newChip)
				}
			}
		}
	case Dragging:
		last := g.pointer.lastPos
		cur := g.pointer.currentPos
		if g.boardWindow.Contains(last) && g.boardWindow.Contains(cur) {
			xx1, yy1 := g.boardWindow.Coords(last)
			xx2, yy2 := g.boardWindow.Coords(cur)
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
}

func (g *Game) updateMaze() {
	switch g.pointer.status {
	case TouchUp:
		if g.pointer.frames > 0 {
			cur := g.pointer.currentPos
			if g.mazeControlsWindow.Contains(cur) {
				xx, yy := g.mazeControlsWindow.Coords(cur)
				g.gameControlSelector.Click(xx, yy)
			}
		}
	}
}

func (g *Game) Draw(screen *ebiten.Image) {
	g.drawBoard(screen)
	g.drawMaze(screen)
}

func (g *Game) drawMaze(screen *ebiten.Image) {
	g.gameControlSelector.Draw(g.mazeControlsWindow.Canvas(screen))
	g.boardController.maze.Draw(g.mazeWindow.Canvas(screen), g.mazeRenderer, float64(g.step)/60, g.count/60)
}

func (g *Game) drawBoard(screen *ebiten.Image) {
	if !g.playing {
		g.chipSelector.Draw(g.boardControlsWindow.Canvas(screen), g.boardRenderer.chips)
	}
	g.board.Draw(g.boardWindow.Canvas(screen), g.boardRenderer)
}

func hSplit(r image.Rectangle, y int) (r1 image.Rectangle, r2 image.Rectangle) {
	r1 = r
	r2 = r
	r1.Max.Y = y
	r2.Min.Y = y
	return
}

func vSplit(r image.Rectangle, x int) (r1 image.Rectangle, r2 image.Rectangle) {
	r1 = r
	r2 = r
	r1.Max.X = x
	r2.Min.X = x
	return
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	g.outsideWidth = outsideWidth
	g.outsideHeight = outsideHeight
	return outsideWidth, outsideHeight
}

func main() {
	var levelFile string
	flag.StringVar(&levelFile, "level", "", "path to r2f level file")
	flag.Parse()
	if levelFile == "" {
		log.Fatal("please provide a level file")
	}
	levelBytes, err := ioutil.ReadFile(levelFile)
	if err != nil {
		log.Fatal("Could not open level file:", err)
	}
	maze, err := MazeFromString(string(levelBytes))
	if err != nil {
		log.Fatalf("Could not create maze: %s", err)
	}
	// 	board, err := CircuitBoardFromString(`
	// |ST -> W? y> TL|
	// | ^    nv     v|
	// |..    MF <- ..|
	// | ^     v      |
	// |.. <- PR      |`)
	// 	if err != nil {
	// 		log.Fatalf("Could not create circuit board: %s", err)
	// 	}
	board := NewCircuitBoard(8, 8)
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
	ebiten.SetWindowSize(1024, 768)
	ebiten.SetWindowTitle("Gobot 2 Flags")
	ebiten.SetWindowResizable(true)
	if err := ebiten.RunGame(game); err != nil {
		log.Fatal(err)
	}
}
