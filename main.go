package main

import (
	"fmt"

	mazemodel "github.com/Drofff/maze-game/maze"
	"github.com/Drofff/maze-game/mazegen"
)

func printCell(cell *mazemodel.Cell) {
	if cell.Walls.Top && cell.Walls.Bottom && cell.Walls.Left && cell.Walls.Right {
		fmt.Print("▢")
		return
	}

	if cell.Walls.Top && cell.Walls.Bottom && cell.Walls.Left {
		fmt.Print("[")
		return
	}
	if cell.Walls.Top && cell.Walls.Bottom && cell.Walls.Right {
		fmt.Print("]")
		return
	}

	if cell.Walls.Top && cell.Walls.Left && cell.Walls.Right {
		fmt.Print("П")
		return
	}
	if cell.Walls.Bottom && cell.Walls.Left && cell.Walls.Right {
		fmt.Print("U")
		return
	}

	if cell.Walls.Top && cell.Walls.Bottom {
		fmt.Print("=")
		return
	}

	if cell.Walls.Top && cell.Walls.Left {
		fmt.Print("Г")
		return
	}

	if cell.Walls.Top && cell.Walls.Right {
		fmt.Print("┓")
		return
	}

	if cell.Walls.Bottom && cell.Walls.Left {
		fmt.Print("L")
		return
	}
	if cell.Walls.Bottom && cell.Walls.Right {
		fmt.Print("┛")
		return
	}

	if cell.Walls.Left && cell.Walls.Right {
		fmt.Print("Н")
		return
	}

	if cell.Walls.Top {
		fmt.Print("‾")
		return
	}
	if cell.Walls.Bottom {
		fmt.Print("_")
		return
	}
	if cell.Walls.Left || cell.Walls.Right {
		fmt.Print("|")
		return
	}

	fmt.Print("err")
}

func main() {
	maze := mazegen.Generate(5, 10)

	var startCell *mazemodel.Cell
	for y := 0; y < len(maze); y++ {
		for x := 0; x < len(maze[y]); x++ {
			if maze[y][x].Role == mazemodel.CellRoleStart {
				startCell = maze[y][x]
			}

			printCell(maze[y][x])
		}
		fmt.Print("\n")
	}

	fmt.Printf("\nstart cell x: %v, y: %v\n", startCell.Loc.ColumnIndex, startCell.Loc.RowIndex)
}
