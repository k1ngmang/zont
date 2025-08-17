package main

import (
	"fmt"
	"strings"
)

type Screen struct {
	matrix *Matrix
}

func NewScreen(matrix *Matrix) *Screen {
	return &Screen{matrix: matrix}
}

func (s *Screen) initScreen(screen [][]rune) {
	for row := 0; row < len(screen); row++ {
		for col := 0; col < len(screen[0])-1; col++ {
			screen[row][col] = ' '
		}
		screen[row][s.matrix.GetCols()] = '\n'
	}
}

func (s *Screen) drawScreen() {
	var buffer strings.Builder

	for row := 0; row < len(s.matrix.screenBuffer[1]); row++ {
		for col := 0; col < len(s.matrix.screenBuffer[1][0]); col++ {
			buffer.WriteRune(s.matrix.screenBuffer[1][row][col])
		}
	}

	fmt.Print("\033[H")
	fmt.Print(buffer.String())
}
