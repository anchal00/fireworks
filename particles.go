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

func getRandomColour() string {
	colours := []string{
		"\033[0m",
		"\033[30m",
		"\033[31m",
		"\033[32m",
		"\033[33m",
		"\033[34m",
		"\033[35m",
		"\033[36m",
		"\033[37m",
	}
	return colours[rand.Intn(len(colours))]
}

func renderFirework(xCoord, yCoord int, window *[ROWS][COLS]string, depth int, writer *bufio.Writer) {
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
	queue := [][]int{{xCoord, yCoord}}

	for depth >= 0 {
		size := len(queue)
		if size == 0 {
			return
		}
		for i := 0; i < size; i++ {
			cur_cell := queue[0]
			cur_cell_x, cur_cell_y := cur_cell[0], cur_cell[1]
			queue = queue[1:]
      window[cur_cell_x][cur_cell_y] = fmt.Sprintf("%v%v", getRandomColour() ,"x")

			for _, offset := range offsets {
				rowOffset, colOffset := offset[0], offset[1]
				newX, newY := rowOffset+cur_cell_x, colOffset+cur_cell_y
				if newX < 0 || newX >= len(window) || newY < 0 || newY >= len(window[0]) {
					continue
				}
				if window[newX][newY] == "x" {
					continue
				}
				queue = append(queue, []int{newX, newY})
			}
		}
		// render here
		write(writer, window)
		depth--
		time.Sleep(200 * time.Millisecond)
	}
}

func write(writer *bufio.Writer, cells *[ROWS][COLS]string) {
	fmt.Fprint(writer, "\033[H")
	for i := 0; i < ROWS; i++ {
		for j := 0; j < COLS; j++ {
			fmt.Fprintf(writer, "%v ", cells[i][j])
		}
		fmt.Fprintln(writer)
	}
	writer.Flush()
}

func render() {
	writer := bufio.NewWriter(os.Stdout)
	for {
		fmt.Fprint(writer, "\033[H") /* Write ANSI escape sequence to move the cursor to top left corner. Although, this
		   sequence is not meant for clearing the screen but it works for me ¯\_(ツ)_/¯ i.e
		   clears the screen. And the actual sequence \033[2J that should clear the screen
		   doesn"t work. */

		writer.Flush() // Flush the clear screen command immediately

		xCoord := int(rand.Uint32() % uint32(ROWS))
		yCoord := int(rand.Uint32() % uint32(COLS))

		cells := [ROWS][COLS]string{}
		for i := 0; i < ROWS; i++ {
			for j := 0; j < COLS; j++ {
				cells[i][j] = " "
			}
		}
		depth := int(rand.Intn(6))
		renderFirework(xCoord, yCoord, &cells, depth, writer)
	}
}
