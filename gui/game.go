package gui

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
	mazemodel "github.com/Drofff/maze-game/maze"
	"github.com/Drofff/maze-game/mazegen"
	"image/color"
)

const (
	mazeWidth  = 50
	mazeHeight = 50

	gameWindowSize = 900
)

func buildPixelMatrix(m [][]*mazemodel.Cell, windowSize int) [][]color.Color {
	pm := make([][]color.Color, windowSize)
	for row := 0; row < windowSize; row++ {
		pm[row] = make([]color.Color, windowSize)
		for column := 0; column < windowSize; column++ {
			pm[row][column] = color.Black
		}
	}

	pixelsPerCell := windowSize / len(m)
	for y := 0; y < len(m); y++ {
		for x := 0; x < len(m[y]); x++ {
			startPixelY := y * pixelsPerCell
			startPixelX := x * pixelsPerCell

			endPixelY := startPixelY + pixelsPerCell - 1
			endPixelX := startPixelX + pixelsPerCell - 1

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

func newGameWindow(a fyne.App) fyne.Window {
	w := a.NewWindow("Maze Game")
	w.Resize(fyne.NewSquareSize(gameWindowSize))
	w.CenterOnScreen()

	m := mazegen.Generate(mazeWidth, mazeHeight)

	var pixelMatrix [][]color.Color
	w.SetContent(canvas.NewRasterWithPixels(func(x, y, w, h int) color.Color {
		if len(pixelMatrix) == 0 {
			pixelMatrix = buildPixelMatrix(m, w)
		}
		return pixelMatrix[y][x]
	}))
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
