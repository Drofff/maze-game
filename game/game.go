package game

import (
	"fmt"
	"math"

	"github.com/Drofff/maze-game/maze"
)

type State int

type Game interface {
	State() State
	PlayerLocation() maze.CellLocation
	MoveTo(loc maze.CellLocation)
}

type game struct {
	state      State
	playerLoc  maze.CellLocation
	mazeScheme [][]*maze.Cell
}

const (
	StateOngoing = iota
	StateWon
)

func NewGame(mazeScheme [][]*maze.Cell, playerLoc maze.CellLocation) Game {
	return &game{
		state:      StateOngoing,
		playerLoc:  playerLoc,
		mazeScheme: mazeScheme,
	}
}

func (g *game) State() State {
	return g.state
}

func (g *game) PlayerLocation() maze.CellLocation {
	return g.playerLoc
}

func (g *game) movesUp(loc maze.CellLocation) bool {
	return loc.RowIndex < g.playerLoc.RowIndex
}

func (g *game) movesDown(loc maze.CellLocation) bool {
	return loc.RowIndex > g.playerLoc.RowIndex
}

func (g *game) movesLeft(loc maze.CellLocation) bool {
	return loc.ColumnIndex < g.playerLoc.ColumnIndex
}

func (g *game) movesRight(loc maze.CellLocation) bool {
	return loc.ColumnIndex > g.playerLoc.ColumnIndex
}

func (g *game) canMoveTo(loc maze.CellLocation) bool {
	if loc.RowIndex < 0 || loc.RowIndex >= len(g.mazeScheme) || loc.ColumnIndex < 0 || loc.ColumnIndex >= len(g.mazeScheme[0]) {
		return false
	}

	dist := math.Abs(float64(loc.ColumnIndex-g.playerLoc.ColumnIndex)) + math.Abs(float64(loc.RowIndex-g.playerLoc.RowIndex))
	if dist != 1 {
		return false
	}

	currCell := g.mazeScheme[g.playerLoc.RowIndex][g.playerLoc.ColumnIndex]
	destCell := g.mazeScheme[loc.RowIndex][loc.ColumnIndex]

	if g.movesUp(loc) {
		return !currCell.Walls.Top && !destCell.Walls.Bottom
	}

	if g.movesDown(loc) {
		return !currCell.Walls.Bottom && !destCell.Walls.Top
	}

	if g.movesLeft(loc) {
		return !currCell.Walls.Left && !destCell.Walls.Right
	}

	if g.movesRight(loc) {
		return !currCell.Walls.Right && !destCell.Walls.Left
	}

	panic(fmt.Sprintf("unexpected condition: pl %v, dl %v", g.playerLoc, loc))
}

func (g *game) MoveTo(loc maze.CellLocation) {
	if !g.canMoveTo(loc) {
		return
	}

	destCell := g.mazeScheme[loc.RowIndex][loc.ColumnIndex]
	if destCell.Role == maze.CellRoleFinish {
		g.state = StateWon
	}

	g.playerLoc = loc
}
