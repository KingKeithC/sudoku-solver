package main

import (
	"bytes"
	"fmt"
)

// Sudoku represents a game of Sudoku.
type Sudoku struct {
	board [][]int
}

// NewPremade returns a premade known solvable board of Sudoku.
func NewPremade() *Sudoku {
	board := [][]int{
		{0, 1, 0, 0, 4, 0, 6, 0, 0},
		{2, 7, 8, 0, 0, 3, 9, 0, 0},
		{3, 0, 4, 0, 0, 1, 2, 0, 0},
		{0, 0, 9, 0, 6, 0, 1, 0, 4},
		{1, 0, 0, 4, 0, 8, 0, 0, 6},
		{6, 0, 7, 0, 3, 0, 5, 0, 0},
		{0, 0, 3, 5, 0, 0, 8, 0, 9},
		{0, 0, 1, 3, 0, 0, 5, 0, 0},
		{0, 0, 6, 0, 2, 0, 0, 1, 0},
	}

	return &Sudoku{
		board: board,
	}
}

// String returns a string representation of s.
func (s *Sudoku) String() string {
	b := bytes.Buffer{}

	for row := 0; row < len(s.board); row++ {

		if row != 0 && (row%3 == 0) {
			b.WriteString("- - - + - - - + - - -\n")
		}

		for col := 0; col < len(s.board[row]); col++ {
			val := s.board[row][col]
			//fmt.Printf("(%d, %d) = %d\n", row, col, val)

			if col != 0 && (col%3 == 0) {
				b.WriteString("| ")
			}

			b.WriteString(fmt.Sprintf("%d ", val))
		}

		b.WriteString("\n")
	}

	return b.String()
}

func main() {
	s := NewPremade()
	fmt.Println(s.String())
}
