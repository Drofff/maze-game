package mazegen

import (
	"math/rand"

	"github.com/golang-collections/collections/stack"
)

type CellWalls struct {
	Top    bool
	Bottom bool
	Left   bool
	Right  bool
}

type CellLocation struct {
	RowIndex    int
	ColumnIndex int
}

type CellRole int

type Cell struct {
	Walls   CellWalls
	Loc     CellLocation
	Role    CellRole
	visited bool
}

const (
	CellRolePath   CellRole = 0
	CellRoleStart  CellRole = 1
	CellRoleFinish CellRole = 2
)

func newEmptyMaze(width, height int) [][]*Cell {
	maze := make([][]*Cell, height)
	for y := 0; y < height; y++ {
		maze[y] = make([]*Cell, width)
		for x := 0; x < width; x++ {
			maze[y][x] = &Cell{
				Walls:   CellWalls{Top: true, Bottom: true, Left: true, Right: true},
				Loc:     CellLocation{RowIndex: y, ColumnIndex: x},
				Role:    CellRolePath,
				visited: false,
			}
		}
	}
	return maze
}

func findUnvisitedNeighbours(posX, posY int, maze [][]*Cell) []*Cell {
	neighbours := []CellLocation{
		{RowIndex: posY, ColumnIndex: posX - 1},
		{RowIndex: posY, ColumnIndex: posX + 1},
		{RowIndex: posY - 1, ColumnIndex: posX},
		{RowIndex: posY + 1, ColumnIndex: posX},
	}

	var unvisitedNeighbours []*Cell
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

func selectRandom(cells []*Cell) *Cell {
	if len(cells) == 1 {
		return cells[0]
	}

	randIndex := rand.Intn(len(cells))
	return cells[randIndex]
}

func removeWallBetween(cellA, cellB *Cell) {
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

func Generate(width, height int) [][]*Cell {
	maze := newEmptyMaze(width, height)
	currPosX, currPosY := rand.Intn(width), 0

	maze[currPosY][currPosX].Role = CellRoleStart
	maze[currPosY][currPosX].visited = true

	maze[height-1][width-1].Role = CellRoleFinish

	visitedCells := stack.New()

	for {
		uns := findUnvisitedNeighbours(currPosX, currPosY, maze)
		if len(uns) == 0 {
			previousCell := visitedCells.Pop()
			if previousCell == nil {
				return maze
			}

			currPosX, currPosY = previousCell.(*Cell).Loc.ColumnIndex, previousCell.(*Cell).Loc.RowIndex
			continue
		}
		visitedCells.Push(maze[currPosY][currPosX])

		nextCell := selectRandom(uns)
		removeWallBetween(maze[currPosY][currPosX], nextCell)
		nextCell.visited = true

		currPosX, currPosY = nextCell.Loc.ColumnIndex, nextCell.Loc.RowIndex
	}
}
