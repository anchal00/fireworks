package main

import (
	"bufio"
	"fmt"
	"math"
	"math/rand"
	"os"
	"os/signal"
	"sync"
	"time"
)

const (
	ROWS = 50
	COLS = 100
)

var outputWriterLock = sync.RWMutex{}

func getRandomColour() string {
	colours := []string{
		"\033[38;5;196m", // Bright Red
		"\033[38;5;208m", // Bright Orange
		"\033[38;5;226m", // Bright Yellow
		"\033[38;5;46m",  // Bright Green
		"\033[38;5;51m",  // Bright Cyan
		"\033[38;5;45m",  // Bright Blue
		"\033[38;5;201m", // Bright Magenta
		"\033[38;5;213m", // Bright Pink
		"\033[38;5;129m", // Bright Purple
		"\033[38;5;154m", // Bright Lime Green
		"\033[38;5;111m", // Bright Sky Blue
		"\033[38;5;51m",  // Bright Turquoise
		"\033[38;5;220m", // Bright Gold
		"\033[38;5;135m", // Bright Violet
		"\033[38;5;203m", // Bright Coral
		"\033[38;5;85m",  // Bright Mint Green
		"\033[38;5;217m", // Bright Light Pink
		"\033[38;5;159m", // Bright Light Blue
		"\033[38;5;216m", // Bright Peach
		"\033[38;5;229m", // Bright Light Yellow
		"\033[38;5;195m", // Bright Light Cyan
		"\033[38;5;156m", // Bright Light Green
		"\033[38;5;99m",  // Bright Indigo
		"\033[38;5;183m", // Bright Lavender
		"\033[38;5;122m", // Bright Aquamarine
	}
	return colours[rand.Intn(len(colours))]
}

func getRandomChar() string {
	markers := []string{
		"~",
		"x",
		"-",
		"*",
		"0",
		"+",
	}
	return markers[rand.Intn(len(markers))]
}

func renderFirework(cx, cy int, window *[ROWS][COLS]string, radius int, writer *bufio.Writer) {
	lastChar := "^"
	for i := ROWS - 1; i > cx; i-- {
		if lastChar == "^" {
			window[i][cy] = "|"
			lastChar = "|"
		} else {
			window[i][cy] = "^"
			lastChar = "^"
		}
		write(writer, window)
		time.Sleep(3 * time.Millisecond)
	}
	// Fade away seq 1
	for i := ROWS - 1; i > cx; i-- {
		window[i][cy] = " "
		write(writer, window)
		time.Sleep(1 * time.Millisecond)
	}
	explosionMarker := getRandomChar()
	for i := 1; i <= radius; i++ {
		colour := getRandomColour()
		for y := 0; y < COLS; y++ {
			for x := 0; x < ROWS; x++ {
				dx := x - cx
				dy := y - cy
				distance := math.Sqrt(float64(dx*dx + dy*dy))
				if distance <= float64(i) && window[x][y] == " " {
					window[x][y] = fmt.Sprintf("%v\033[1m%v", colour, explosionMarker) // \033[1m makes text bold
				}
			}
		}
		write(writer, window)
		time.Sleep(100 * time.Millisecond)
	}
	// Fade away seq 2
	for i := 1; i <= radius; i++ {
		for y := COLS - 1; y >= 0; y-- {
			for x := ROWS - 1; x >= 0; x-- {
				dx := x - cx
				dy := y - cy
				distance := math.Sqrt(float64(dx*dx + dy*dy))
				if distance <= float64(i) {
					window[x][y] = " "
				}
			}
		}
		write(writer, window)
		time.Sleep(120 * time.Millisecond)
	}
}

func write(writer *bufio.Writer, cells *[ROWS][COLS]string) {
	outputWriterLock.Lock()
	for i := 0; i < ROWS; i++ {
		for j := 0; j < COLS; j++ {
			fmt.Fprintf(writer, "%v ", cells[i][j])
		}
		fmt.Fprintln(writer)
	}
	fmt.Fprint(writer, "\033[H\033[0m") // "\033[0m" -> Reset text formatting
	writer.Flush()                      // Flush the clear screen command immediately
	outputWriterLock.Unlock()
}

func create(cells *[ROWS][COLS]string, writer *bufio.Writer, wg *sync.WaitGroup) {
	defer wg.Done()
	xCoord := int(rand.Uint32() % uint32(ROWS))
	yCoord := int(rand.Uint32() % uint32(COLS))

	// center the coordinates
	xCoord = max(xCoord, 20)
	xCoord = min(xCoord, 30)
	yCoord = max(yCoord, 20)
	yCoord = min(yCoord, 80)
	for i := 0; i < ROWS; i++ {
		for j := 0; j < COLS; j++ {
			cells[i][j] = " "
		}
	}
	depth := int(rand.Intn(6)) + 2
	renderFirework(xCoord, yCoord, cells, depth, writer)
}

func render() {
	cells := [ROWS][COLS]string{}
	writer := bufio.NewWriter(os.Stdout)
	channel := make(chan os.Signal, 1)
	signal.Notify(channel, os.Interrupt)
	go func() {
		<-channel
		fmt.Fprintln(writer, "\033[2J\033[H\033[0m\n\nExiting.....Goodbye !\n")
		writer.Flush() // Flush the clear screen command immediately
		os.Exit(0)
	}()

	for {
		fmt.Fprint(writer, "\033[H\033[0m") /* Write ANSI escape sequence to move the cursor to top left corner. Although, this
		   sequence is not meant for clearing the screen but it works for me ¯\_(ツ)_/¯ i.e
		   clears the screen. And the actual sequence \033[2J that should clear the screen
		   doesn"t work. */
		writer.Flush() // Flush the clear screen command immediately
		fireworks := rand.Intn(5) + 1
		wg := sync.WaitGroup{}
		for i := 0; i <= fireworks; i++ {
			wg.Add(1)
			go create(&cells, writer, &wg)
		}
		wg.Wait()
	}
}
