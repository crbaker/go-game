package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"os/exec"
	"time"
)

type grid struct {
	grid [][]int
}

func (g *grid) surroundingCells(row int, col int) []int {
	surrounding := []int{}

	for r := row - 1; r <= row+1; r++ {
		if (r < 0) || (r >= len(g.grid)) {
			continue
		}

		record := g.grid[r]
		for c := col - 1; c <= col+1; c++ {
			if (c < 0) || (c >= len(record)) {
				continue
			}
			if (r == row) && (c == col) {
				continue
			}
			surrounding = append(surrounding, record[c])
		}
	}

	return surrounding
}

func (g *grid) liveNeighbours(row int, col int) int {
	surrounding := g.surroundingCells(row, col)
	return sum(surrounding)
}

func (g *grid) isCellAlive(row int, col int) bool {
	return g.grid[row][col] == 1
}

func (g *grid) determineFate(row int, col int) int {
	count := g.liveNeighbours(row, col)

	if g.isCellAlive(row, col) {

		if count < 2 {
			return 0
		} else if (count == 2) || (count == 3) {
			return 1
		} else if count > 3 {
			return 0
		}

	} else {
		if count == 3 {
			return 1
		}
	}

	return 0
}

func (g *grid) tick() {

	newGrid := make([][]int, len(g.grid))
	copy(newGrid, g.grid)

	for r, row := range g.grid {
		for c := range row {
			newGrid[r][c] = g.determineFate(r, c)
		}
	}
	// g.grid = newGrid
	copy(g.grid, newGrid)
}

func sum(slice []int) int {
	sum := 0
	for _, v := range slice {
		sum += v
	}
	return sum
}

func (g *grid) printGrid() {
	for _, row := range g.grid {
		for _, col := range row {
			if col == 1 {
				fmt.Print("#")
			} else {
				fmt.Print(" ")
			}
		}
		fmt.Println("")
	}
}

func newGrid(filePath string) *grid {
	lines := linesInFile(filePath)
	g := [][]int{}
	if len(*lines) == 0 {
		log.Fatal("the supplied file has no seed data")
	}

	for _, l := range *lines {
		row := []int{}
		for _, c := range l {
			state := 0
			if string(c) == "#" {
				state = 1
			}
			row = append(row, state)
		}
		g = append(g, row)
	}
	return &grid{grid: g}
}

func linesInFile(filePath string) *[]string {
	file, err := os.Open(filePath)

	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	lines := []string{}
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {

		lines = append(lines, scanner.Text())
	}

	return &lines
}

func main() {

	filePath := os.Args[1]

	g := newGrid(filePath)

	for r := 0; r <= 10000; r++ {
		g.printGrid()
		g.tick()
		time.Sleep(100 * time.Millisecond)
		g.printGrid()
		clearConsole()
	}
}

func clearConsole() {
	cmd := exec.Command("clear")
	cmd.Stdout = os.Stdout
	cmd.Run()
}
