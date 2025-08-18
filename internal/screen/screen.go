package screen

import (
	"fmt"
	"strings"
	"zontengine/internal/matrix"
)

type Screen struct {
	matrix *matrix.Matrix
}

func NewScreen(matrix *matrix.Matrix) *Screen {
	return &Screen{matrix: matrix}
}

func (s *Screen) InitScreen(screen [][]rune) {
	for row := 0; row < len(screen); row++ {
		for col := 0; col < len(screen[0])-1; col++ {
			screen[row][col] = ' '
		}
		screen[row][s.matrix.GetCols()] = '\n'
	}
}

func (s *Screen) DrawScreen() {
	var buffer strings.Builder

	for row := 0; row < len(s.matrix.ScreenBuffer[1]); row++ {
		for col := 0; col < len(s.matrix.ScreenBuffer[1][0]); col++ {
			buffer.WriteRune(s.matrix.ScreenBuffer[1][row][col])
		}
	}

	fmt.Print("\033[H")
	fmt.Print(buffer.String())
}
