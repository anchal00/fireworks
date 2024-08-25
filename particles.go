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

func render() {
	writer := bufio.NewWriter(os.Stdout)

	for {
		fmt.Fprint(writer, "\033[H") /* Write ANSI escape sequence to move the cursor to top left corner. Although, this
		   sequence is not meant for clearing the screen but it works for me ¯\_(ツ)_/¯ i.e
		   clears the screen. And the actual sequence \033[2J that should clear the screen
		   doesn't work. */

		writer.Flush() // Flush the clear screen command immediately

		x_coord := int(rand.Uint32() % uint32(ROWS))
		y_coord := int(rand.Uint32() % uint32(COLS))

		cells := [ROWS][COLS]rune{}

		for i := 0; i < ROWS; i++ {
			for j := 0; j < COLS; j++ {
				cells[i][j] = '-'
			}
		}
		cells[x_coord][y_coord] = 'O'
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
