package main

import (
	"bytes"
	"context"
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"

	"github.com/chromedp/cdproto/cdp"
	"github.com/chromedp/chromedp"
	"golang.org/x/net/html"
)

const (
	pause              = time.Second * 3
	timeout            = time.Second * 60
	sudokuURL          = "https://nine.websudoku.com/?level=4"
	puzzleGridSelector = `#puzzle_grid`
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

// NewFromHTML returns an instance of Sudoku with the board initially set from the contents of an HTML table grid, or an error if it could not parse the HTML.
func NewFromHTML(h *string) (*Sudoku, error) {
	node, err := html.Parse(strings.NewReader(*h))
	if err != nil {
		return nil, err
	}

	board := [][]int{
		{0, 0, 0, 0, 0, 0, 0, 0, 0},
		{0, 0, 0, 0, 0, 0, 0, 0, 0},
		{0, 0, 0, 0, 0, 0, 0, 0, 0},
		{0, 0, 0, 0, 0, 0, 0, 0, 0},
		{0, 0, 0, 0, 0, 0, 0, 0, 0},
		{0, 0, 0, 0, 0, 0, 0, 0, 0},
		{0, 0, 0, 0, 0, 0, 0, 0, 0},
		{0, 0, 0, 0, 0, 0, 0, 0, 0},
		{0, 0, 0, 0, 0, 0, 0, 0, 0},
	}

	tbody := node.FirstChild.LastChild.FirstChild.FirstChild // Selects the html/body/table node
	for tr := tbody.FirstChild; tr != nil; tr = tr.NextSibling {
		for td := tr.FirstChild; td != nil; td = td.NextSibling {
			input := td.FirstChild

			// for each `input` tag
			if input.Type == html.ElementNode && input.Data == "input" {
				var (
					id    string
					value int
				)

				// Iterate over the attributes, and select the id, and value
				for _, attr := range input.Attr {
					switch attr.Key {
					case "id":
						id = attr.Val
					case "value":
						value, err = strconv.Atoi(attr.Val)
						if err != nil {
							value = 0
						}
					}
				}

				// Determine the row, and column from the id
				id = id[1:]                      // Strip off first character
				row, err := strconv.Atoi(id[1:]) // Row is second character
				if err != nil {
					return nil, err
				}
				col, err := strconv.Atoi(id[:1]) // Col is first character
				if err != nil {
					return nil, err
				}

				// Set the position on the board to the value
				board[row][col] = value
			}
		}
	}

	return &Sudoku{board}, nil
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

// solveSudoku is intended to be used with
func solveSudoku(ctx context.Context) error {
	var (
		loc            string
		title          string
		tableOuterHTML string
	)

	// Create a new context with a 60 second timeout
	newCtx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	// Connect to the page, and retrieve the puzzle's HTML
	log.Printf("Starting Sudoku Solver with a %1.f second timeout...", timeout.Seconds())
	if err := chromedp.Run(newCtx,
		chromedp.Navigate(sudokuURL),
		chromedp.Location(&loc),
		chromedp.Title(&title),
		chromedp.OuterHTML(puzzleGridSelector, &tableOuterHTML, chromedp.NodeVisible),
	); err != nil {
		return err
	}

	// Parse the HTML into a Sudoku object
	log.Printf("Running Sudoku Solver at [%s] (%s)...", title, loc)
	s, err := NewFromHTML(&tableOuterHTML)
	if err != nil {
		return err
	}

	// Solve the puzzle and print it
	solvable := s.Solve()
	fmt.Println("Solved: ", solvable, "-----------------")
	fmt.Println(s)
	fmt.Println("-----------------------------")

	// Enter the Solution into the page
	log.Println("Entering Solution into page...")
	for row := 0; row < 9; row++ {
		for col := 0; col < 9; col++ {
			var (
				id    = fmt.Sprintf("#f%d%d", col, row) // Determine the ID of the cell for the given position
				value = fmt.Sprintf("%d\n", s.board[row][col])
				nodes = &[]*cdp.Node{}
			)

			// Find the nodes matching the selector
			if chromedp.Run(newCtx, chromedp.Nodes(id, nodes, chromedp.AtLeast(1))); err != nil {
				return err
			}

			// Bail if there were no nodes
			if nodes == nil {
				return fmt.Errorf("no nodes retireved with selector id = \"%s\"", id)
			}

			// Enter the value and wait a bit
			if err := chromedp.Run(newCtx,
				chromedp.KeyEventNode((*nodes)[0], value),
			); err != nil {
				return err
			}

		}
	}

	return nil
}

func main() {
	// Disable the headless option
	opts := append(chromedp.DefaultExecAllocatorOptions[:], chromedp.Flag("headless", false))

	// Create a ExecAllocator that uses the custom options
	allocCtx, cancel := chromedp.NewExecAllocator(context.Background(), opts...)
	defer cancel()

	// Create a new context with the ExecAllocator
	ctx, cancel := chromedp.NewContext(allocCtx)
	defer cancel()

	for {
		err := chromedp.Run(ctx, chromedp.ActionFunc(solveSudoku))
		if err != nil {
			log.Fatal(err)
		}
	}
}
