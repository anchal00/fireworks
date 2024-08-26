package main

import (
	"bufio"
	"fmt"
	"math"
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
		"\033[30m", // Black text
		"\033[31m", // Red text
		"\033[32m", // Green text
		"\033[33m", // Yellow text
		"\033[34m", // Blue text
		"\033[35m", // Magenta text
		"\033[36m", // Cyan text
		"\033[37m", // White text
	}
	return colours[rand.Intn(len(colours))]
}

func renderFirework(cx, cy int, window *[ROWS][COLS]string, radius int, writer *bufio.Writer) {
	for i := 1; i <= radius; i++ {
		colour := getRandomColour()
		for y := 0; y < COLS; y++ {
			for x := 0; x < ROWS; x++ {
				dx := x - cx
				dy := y - cy
				distance := math.Sqrt(float64(dx*dx + dy*dy))
				if distance <= float64(i) {
					window[x][y] = fmt.Sprintf("%v\033[1m%v", colour, "+") // \033[1m makes text bold
				}
			}
		}
		write(writer, window)
		time.Sleep(200 * time.Millisecond)
	}
}

func write(writer *bufio.Writer, cells *[ROWS][COLS]string) {
	for i := 0; i < ROWS; i++ {
		for j := 0; j < COLS; j++ {
			fmt.Fprintf(writer, "%v ", cells[i][j])
		}
		fmt.Fprintln(writer)
	}
	fmt.Fprint(writer, "\033[H\033[0m") // "\033[0m" -> Reset text formatting
	writer.Flush()                      // Flush the clear screen command immediately
}

func render() {
	writer := bufio.NewWriter(os.Stdout)

	for {
		fmt.Fprint(writer, "\033[H\033[0m") /* Write ANSI escape sequence to move the cursor to top left corner. Although, this
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
