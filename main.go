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
		{1, 2, 0, 0, 0, 0, 0, 4, 0},
		{0, 0, 8, 1, 0, 0, 0, 0, 5},
		{0, 4, 0, 2, 0, 9, 1, 6, 3},
		{0, 0, 2, 6, 5, 0, 4, 0, 1},
		{0, 0, 0, 0, 0, 0, 0, 0, 0},
		{8, 0, 6, 0, 7, 1, 5, 0, 0},
		{5, 9, 4, 7, 0, 3, 0, 8, 0},
		{6, 0, 0, 0, 0, 5, 3, 0, 0},
		{0, 1, 0, 0, 0, 0, 0, 5, 9},
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

			// Append the value of the cell, or a blank space if it is empty (val==0)
			if val == 0 {
				b.WriteString("  ")
			} else {
				b.WriteString(fmt.Sprintf("%d ", val))
			}
		}

		if row != 8 {
			b.WriteString("\n")
		}
	}

	return b.String()
}

// isValid takes a board, a position, and a guess and returns whether it is valid.
// The guess should not already be in the board.
func (s *Sudoku) isValid(row, col, guess int) bool {
	// Search each col in the row for the guess.
	for c := 0; c < len(s.board[row]); c++ {
		if s.board[row][c] == guess {
			return false
		}
	}

	// Search each row in the col for the guess.
	for r := 0; r < len(s.board); r++ {
		if s.board[r][col] == guess {
			return false
		}
	}

	// Search the subregion for the guess.
	subregionRow := row - row%3
	subregionCol := col - col%3
	for r := subregionRow; r < subregionRow+3; r++ {
		for c := subregionCol; c < subregionCol+3; c++ {
			if s.board[r][c] == guess {
				return false
			}
		}
	}

	// no guess found, this guess is valid
	return true
}

// Solve replaces missing numbers in the Sudoku board.
func (s *Sudoku) Solve() bool {
	for row := 0; row < 9; row++ {
		for col := 0; col < 9; col++ {
			if s.board[row][col] == 0 {
				for guess := 1; guess <= 9; guess++ {
					if s.isValid(row, col, guess) {
						s.board[row][col] = guess
						if s.Solve() {
							return true
						}
						s.board[row][col] = 0
					}
				}
				return false
			}
		}
	}
	return true
}

func main() {
	s := NewPremade()

	fmt.Println("Initial Board:")
	fmt.Println("----------------------------")
	fmt.Println(s)
	fmt.Printf("----------------------------\n\n")

	fmt.Printf("Solvable: %t\n", s.Solve())
	fmt.Println("----------------------------")
	fmt.Println(s)
	fmt.Printf("----------------------------\n\n")
}
