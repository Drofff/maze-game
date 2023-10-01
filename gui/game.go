package gui

import (
	"image/color"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
	"github.com/Drofff/maze-game/game"
	mazemodel "github.com/Drofff/maze-game/maze"
	"github.com/Drofff/maze-game/mazegen"
)

const (
	mazeWidth  = 50
	mazeHeight = 50

	gameWindowSize = 900

	keyScanCodeUp    = 13 // 'W' key mac os
	keyScanCodeDown  = 1  // 'S' key mac os
	keyScanCodeLeft  = 0  // 'A' key mac os
	keyScanCodeRight = 2  // 'D' key mac os
)

func newWinnerWindow(a fyne.App) fyne.Window {
	w := a.NewWindow("Congrats!")
	img := canvas.NewImageFromFile("img/you-win.jpeg")
	img.FillMode = canvas.ImageFillOriginal
	w.SetContent(container.NewStack(img))
	w.CenterOnScreen()
	return w
}

func buildPixelMatrix(m [][]*mazemodel.Cell, windowSize int) [][]color.Color {
	pm := make([][]color.Color, windowSize)
	for row := 0; row < windowSize; row++ {
		pm[row] = make([]color.Color, windowSize)
		for column := 0; column < windowSize; column++ {
			pm[row][column] = color.Black
		}
	}

	pixelsPerCell := float32(windowSize) / float32(len(m))
	for y := 0; y < len(m); y++ {
		for x := 0; x < len(m[y]); x++ {
			startPixelY := int(float32(y) * pixelsPerCell)
			startPixelX := int(float32(x) * pixelsPerCell)

			endPixelY := int(float32(startPixelY) + pixelsPerCell)
			endPixelX := int(float32(startPixelX) + pixelsPerCell)

			cell := m[y][x]

			var cellColor color.Color = color.White
			switch cell.Role {
			case mazemodel.CellRoleStart:
				cellColor = color.Black
			case mazemodel.CellRoleFinish:
				cellColor = color.RGBA{R: 135, G: 211, B: 124, A: 255}
			}

			if cell.Walls.Top {
				for wallX := startPixelX; wallX <= endPixelX; wallX++ {
					pm[startPixelY][wallX] = cellColor
				}
			}

			if cell.Walls.Bottom {
				for wallX := startPixelX; wallX <= endPixelX; wallX++ {
					pm[endPixelY][wallX] = cellColor
				}
			}

			if cell.Walls.Left {
				for wallY := startPixelY; wallY <= endPixelY; wallY++ {
					pm[wallY][startPixelX] = cellColor
				}
			}

			if cell.Walls.Right {
				for wallY := startPixelY; wallY <= endPixelY; wallY++ {
					pm[wallY][endPixelX] = cellColor
				}
			}
		}
	}

	return pm
}

func calcPlayerScreenPosition(g game.Game, pixelsInCell float32) fyne.Position {
	playerLoc := g.PlayerLocation()
	return fyne.NewPos(float32(playerLoc.ColumnIndex)*pixelsInCell, float32(playerLoc.RowIndex)*pixelsInCell)
}

func newGameWindow(a fyne.App) fyne.Window {
	w := a.NewWindow("Maze Game")
	w.Resize(fyne.NewSquareSize(gameWindowSize))
	w.CenterOnScreen()

	m := mazegen.Generate(mazeWidth, mazeHeight)

	var pixelMatrix [][]color.Color
	mazeImage := canvas.NewRasterWithPixels(func(x, y, w, h int) color.Color {
		if len(pixelMatrix) == 0 {
			pixelMatrix = buildPixelMatrix(m, w)
		}
		return pixelMatrix[y][x]
	})

	g := game.NewGame(m)

	pixelsInCell := float32(gameWindowSize) / float32(len(m))

	playerMarker := canvas.NewCircle(color.RGBA{R: 255, G: 252, B: 127, A: 255})
	playerMarker.Resize(fyne.NewSquareSize(pixelsInCell / 2))
	playerMarker.Move(calcPlayerScreenPosition(g, pixelsInCell))

	w.SetContent(container.NewStack(mazeImage, container.NewWithoutLayout(playerMarker)))

	w.Canvas().SetOnTypedKey(func(e *fyne.KeyEvent) {
		playerLoc := g.PlayerLocation()
		switch e.Physical.ScanCode {
		case keyScanCodeUp:
			g.MoveTo(mazemodel.CellLocation{RowIndex: playerLoc.RowIndex - 1, ColumnIndex: playerLoc.ColumnIndex})
		case keyScanCodeDown:
			g.MoveTo(mazemodel.CellLocation{RowIndex: playerLoc.RowIndex + 1, ColumnIndex: playerLoc.ColumnIndex})
		case keyScanCodeLeft:
			g.MoveTo(mazemodel.CellLocation{RowIndex: playerLoc.RowIndex, ColumnIndex: playerLoc.ColumnIndex - 1})
		case keyScanCodeRight:
			g.MoveTo(mazemodel.CellLocation{RowIndex: playerLoc.RowIndex, ColumnIndex: playerLoc.ColumnIndex + 1})
		}

		if g.PlayerLocation().RowIndex != playerLoc.RowIndex || g.PlayerLocation().ColumnIndex != playerLoc.ColumnIndex {
			playerMarker.Move(calcPlayerScreenPosition(g, pixelsInCell))
		}

		if g.State() == game.StateWon {
			newWinnerWindow(a).Show()
			w.Close()
		}
	})

	return w
}

func newMenuWindow(a fyne.App) fyne.Window {
	w := a.NewWindow("MAZE - Game Menu")
	w.Resize(fyne.NewSize(400, 300))
	w.CenterOnScreen()

	logo := canvas.NewImageFromFile("img/maze-title.png")
	logo.FillMode = canvas.ImageFillOriginal
	sb := widget.NewButton("Start", func() {
		newGameWindow(a).Show()
		w.Close()
	})

	c := container.NewVBox(widget.NewLabel(""), logo, widget.NewLabel(""), sb)
	w.SetContent(c)
	return w
}

func StartGame() {
	a := app.New()
	w := newMenuWindow(a)
	w.ShowAndRun()
}
