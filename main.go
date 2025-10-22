package main

import (
	"bufio"
	"fmt"
	"math/rand"
	"os"
	"strconv"
	"strings"
)

type Direction int

const (
	Up Direction = iota
	Right
	Down
	Left
)

type Point struct {
	R int
	C int
}

func dirDelta(d Direction) (int, int) {
	switch d {
	case Up:
		return -1, 0
	case Right:
		return 0, 1
	case Down:
		return 1, 0
	case Left:
		return 0, -1
	default:
		return 0, 0
	}
}

var dirRunes = map[Direction]rune{
	Up: '^', Right: '>', Down: 'v', Left: '<',
}

func parseCommandToDir(ch byte) (Direction, bool) {
	switch ch {
	case 'w', 'W':
		return Up, true
	case 'd', 'D':
		return Right, true
	case 's', 'S':
		return Down, true
	case 'a', 'A':
		return Left, true
	default:
		return Up, false
	}
}

type Game struct {
	Grid   [][]rune
	Player Point
	Facing Direction
	Steps  int
	Won    bool
	Rows   int
	Cols   int
}

func main() {
	fmt.Println("Welcome to the Grid!")
	gridSize, density := readDimensions()

	game, err := newGame(gridSize, density)
	if err != nil {
		fmt.Println("Error starting game:", err)
		return
	}

	reader := bufio.NewReader(os.Stdin)

	// The game loop is now much simpler
	for {
		clearScreen()
		game.render()
		fmt.Print("Move (W/A/S/D), Q to quit: ")

		line, _ := reader.ReadString('\n')
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}

		ch := line[0]
		if ch == 'q' || ch == 'Q' {
			fmt.Println("Goodbye!")
			return
		}

		if dir, ok := parseCommandToDir(ch); ok {
			game.TryMove(dir)
			if game.Won {
				clearScreen()
				game.render()
				fmt.Printf("You won! Total steps: %d\n", game.Steps)
				return
			}
		}
	}
}

func newGame(size int, density float64) (*Game, error) {
	if size <= 0 {
		return nil, fmt.Errorf("invalid grid size")
	}

	grid := generateGrid(size, density)
	pos, facing, err := randomPlayerPlacement(grid)
	if err != nil {
		return nil, err
	}

	return &Game{
		Grid:   grid,
		Player: pos,
		Facing: facing,
		Steps:  0,
		Won:    false,
		Rows:   size,
		Cols:   size,
	}, nil
}

func readDimensions() (int, float64) {
	reader := bufio.NewReader(os.Stdin)
	fmt.Print("enter grid dimension (e.g., 10 for 10x10): ")
	dimension, _ := reader.ReadString('\n')
	dimension = strings.TrimSpace(dimension)

	n, err := strconv.Atoi(dimension)
	if err != nil {
		fmt.Println("Invalid input. Please enter an integer like 10.")
		return 0, 0
	}

	fmt.Print("enter the density (e.g. .2 for 20%): ")
	density, _ := reader.ReadString('\n')
	density = strings.TrimSpace(density)
	percent, err := strconv.ParseFloat(density, 64)
	if err != nil {
		fmt.Println("invalid input please enter a decimal point less than 1 like .2")
		return 0, 0
	}
	return n, percent
}

func generateGrid(size int, density float64) [][]rune {
	if density < 0 {
		density = 0
	}
	if density > 1 {
		density = 1
	}
	grid := make([][]rune, size)
	for i := 0; i < size; i++ {
		row := make([]rune, size)
		for j := 0; j < size; j++ {
			if rand.Float64() < density {
				row[j] = '#'
			} else {
				row[j] = '.'
			}
		}
		grid[i] = row
	}
	return grid
}

func emptyCells(grid [][]rune) []Point {
	var cells []Point
	for r, row := range grid {
		for c, ch := range row {
			if ch == '.' {
				cells = append(cells, Point{R: r, C: c})
			}
		}
	}
	return cells
}

func randomPlayerPlacement(grid [][]rune) (Point, Direction, error) {
	empties := emptyCells(grid)
	if len(empties) == 0 {
		return Point{}, 0, fmt.Errorf("no empty cells generated")
	}
	pos := empties[rand.Intn(len(empties))]
	dir := Direction(rand.Intn(4))
	return pos, dir, nil
}

// func renderWithPlayer(grid [][]rune, pos Point, facing Direction) {
// 	for r := range grid {
// 		line := make([]rune, len(grid[r]))
// 		copy(line, grid[r])
// 		if r == pos.R {
// 			line[pos.C] = dirRunes[facing]
// 		}
// 		fmt.Println(string(line))
// 	}
// }

func (g *Game) render() {
	for r := range g.Grid {
		line := make([]rune, g.Cols)
		copy(line, g.Grid[r])
		if r == g.Player.R {
			line[g.Player.C] = dirRunes[g.Facing]
		}
		fmt.Println(string(line))
	}
	fmt.Printf("Steps: %d\n", g.Steps)
}

// func canStep(grid [][]rune, pos Point, dir Direction) (Point, bool, bool) {
// 	dr, dc := dirDelta(dir)
// 	nr, nc := pos.R+dr, pos.C+dc
// 	rows := len(grid)
// 	cols := len(grid[0])

// 	if nr < 0 || nr >= rows || nc < 0 || nc >= cols {
// 		return Point{nr, nc}, true, false
// 	}
// 	if grid[nr][nc] == '#' {
// 		return Point{nr, nc}, false, true
// 	}
// 	return Point{nr, nc}, false, false
// }

// func tryMove(
// 	grid [][]rune,
// 	pos Point,
// 	facing Direction,
// 	dir Direction,
// 	steps int,
// ) (Point, Direction, int, bool) {
// 	next, offMap, blocked := canStep(grid, pos, dir)
// 	if offMap {
// 		return pos, dir, steps + 1, true
// 	}
// 	if blocked {
// 		return pos, dir, steps, false
// 	}
// 	return next, dir, steps + 1, false
// }

func (g *Game) TryMove(dir Direction) {
	dr, dc := dirDelta(dir)
	nr, nc := g.Player.R+dr, g.Player.C+dc

	g.Facing = dir // Always update facing

	// Check for win condition (off-map)
	if nr < 0 || nr >= g.Rows || nc < 0 || nc >= g.Cols {
		g.Steps++
		g.Won = true
		return
	}

	// Check for wall
	if g.Grid[nr][nc] == '#' {
		// Blocked, do nothing more
		return
	}

	// Successful move
	g.Player.R = nr
	g.Player.C = nc
	g.Steps++
}

func clearScreen() {
	// \x1b[2J clears the entire screen
	// \x1b[H moves the cursor to the home position (top-left)
	fmt.Print("\x1b[2J\x1b[H")
}
