package main

import (
	"bufio"
	"fmt"
	"math/rand"
	"os"
	"time"
)

const (
	ROWS = 50
	COLS = 100
)

func placeFirework(xCoord, yCoord int, window *[ROWS][COLS]rune) {
	window[xCoord][yCoord] = 'O'
	// 8 adjacent cells to x,y pos
	offsets := [8][2]int{
		{0, 1},
		{0, -1},
		{-1, 0},
		{1, 0},
		{-1, -1},
		{-1, 1},
		{1, -1},
		{1, 1},
	}
	for _, offset := range offsets {
		rowOffset, colOffset := offset[0], offset[1]
		newX := int(xCoord) + rowOffset
		newY := int(yCoord) + colOffset
		if newX < 0 || newX >= len(window) || newY < 0 || newY >= len(window[0]) {
			continue
		}
		window[newX][newY] = 'A'
	}
}

func render() {
	writer := bufio.NewWriter(os.Stdout)

	for {
		fmt.Fprint(writer, "\033[H") /* Write ANSI escape sequence to move the cursor to top left corner. Although, this
		   sequence is not meant for clearing the screen but it works for me ¯\_(ツ)_/¯ i.e
		   clears the screen. And the actual sequence \033[2J that should clear the screen
		   doesn't work. */

		writer.Flush() // Flush the clear screen command immediately

		xCoord := int(rand.Uint32() % uint32(ROWS))
		yCoord := int(rand.Uint32() % uint32(COLS))

		cells := [ROWS][COLS]rune{}

		for i := 0; i < ROWS; i++ {
			for j := 0; j < COLS; j++ {
				cells[i][j] = '-'
			}
		}
		placeFirework(xCoord, yCoord, &cells)
		for i := 0; i < ROWS; i++ {
			for j := 0; j < COLS; j++ {
				fmt.Fprintf(writer, "%c ", cells[i][j])
			}
			fmt.Fprintln(writer)
		}
		writer.Flush()

		time.Sleep(1 * time.Second)
	}
}
