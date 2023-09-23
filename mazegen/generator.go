package mazegen

import (
	"math/rand"

	mazemodel "github.com/Drofff/maze-game/maze"
	"github.com/golang-collections/collections/stack"
)

type cellWithMeta struct {
	mazemodel.Cell

	visited bool
}

func newEmptyMaze(width, height int) [][]*cellWithMeta {
	maze := make([][]*cellWithMeta, height)
	for y := 0; y < height; y++ {
		maze[y] = make([]*cellWithMeta, width)
		for x := 0; x < width; x++ {
			maze[y][x] = &cellWithMeta{
				Cell: mazemodel.Cell{
					Walls: mazemodel.CellWalls{Top: true, Bottom: true, Left: true, Right: true},
					Loc:   mazemodel.CellLocation{RowIndex: y, ColumnIndex: x},
					Role:  mazemodel.CellRolePath,
				},
				visited: false,
			}
		}
	}
	return maze
}

func findUnvisitedNeighbours(posX, posY int, maze [][]*cellWithMeta) []*cellWithMeta {
	neighbours := []mazemodel.CellLocation{
		{RowIndex: posY, ColumnIndex: posX - 1},
		{RowIndex: posY, ColumnIndex: posX + 1},
		{RowIndex: posY - 1, ColumnIndex: posX},
		{RowIndex: posY + 1, ColumnIndex: posX},
	}

	var unvisitedNeighbours []*cellWithMeta
	for _, neighbour := range neighbours {
		if neighbour.RowIndex < 0 || neighbour.RowIndex >= len(maze) ||
			neighbour.ColumnIndex < 0 || neighbour.ColumnIndex >= len(maze[neighbour.RowIndex]) {
			continue
		}

		cell := maze[neighbour.RowIndex][neighbour.ColumnIndex]
		if !cell.visited {
			unvisitedNeighbours = append(unvisitedNeighbours, cell)
		}
	}

	return unvisitedNeighbours
}

func selectRandom(cells []*cellWithMeta) *cellWithMeta {
	if len(cells) == 1 {
		return cells[0]
	}

	randIndex := rand.Intn(len(cells))
	return cells[randIndex]
}

func removeWallBetween(cellA, cellB *cellWithMeta) {
	if cellA.Loc.RowIndex < cellB.Loc.RowIndex {
		cellA.Walls.Bottom = false
		cellB.Walls.Top = false
		return
	}

	if cellA.Loc.RowIndex > cellB.Loc.RowIndex {
		cellA.Walls.Top = false
		cellB.Walls.Bottom = false
		return
	}

	if cellA.Loc.ColumnIndex < cellB.Loc.ColumnIndex {
		cellA.Walls.Right = false
		cellB.Walls.Left = false
		return
	}

	if cellA.Loc.ColumnIndex > cellB.Loc.ColumnIndex {
		cellA.Walls.Left = false
		cellB.Walls.Right = false
		return
	}
}

func fromCellsWithMeta(cellsWithMeta [][]*cellWithMeta) [][]*mazemodel.Cell {
	cells := make([][]*mazemodel.Cell, len(cellsWithMeta))
	for row := 0; row < len(cellsWithMeta); row++ {
		cells[row] = make([]*mazemodel.Cell, len(cellsWithMeta[row]))

		for column := 0; column < len(cellsWithMeta[row]); column++ {
			cells[row][column] = &cellsWithMeta[row][column].Cell
		}
	}
	return cells
}

func Generate(width, height int) [][]*mazemodel.Cell {
	maze := newEmptyMaze(width, height)
	currPosX, currPosY := rand.Intn(width), 0

	maze[currPosY][currPosX].Role = mazemodel.CellRoleStart
	maze[currPosY][currPosX].visited = true

	maze[height-1][width-1].Role = mazemodel.CellRoleFinish

	visitedCells := stack.New()

	for {
		uns := findUnvisitedNeighbours(currPosX, currPosY, maze)
		if len(uns) == 0 {
			previousCell := visitedCells.Pop()
			if previousCell == nil {
				return fromCellsWithMeta(maze)
			}

			currPosX, currPosY = previousCell.(*cellWithMeta).Loc.ColumnIndex, previousCell.(*cellWithMeta).Loc.RowIndex
			continue
		}
		visitedCells.Push(maze[currPosY][currPosX])

		nextCell := selectRandom(uns)
		removeWallBetween(maze[currPosY][currPosX], nextCell)
		nextCell.visited = true

		currPosX, currPosY = nextCell.Loc.ColumnIndex, nextCell.Loc.RowIndex
	}
}
