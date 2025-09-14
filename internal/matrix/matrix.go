package matrix

/**
 * Provides matrix operations and screen buffer management for 3D rendering.
 * Handles mathematical transformations, z-sorting of triangles, and screen buffering.
 *
 * @param cols          number of columns in the terminal screen
 * @param rows          number of rows in the terminal screen
 * @param angle         current rotation angle for 3D transformations
 * @param ScreenBuffer  double buffer system for smooth rendering [0]=current, [1]=previous
 *
 * The matrix supports:
 * - Matrix multiplication for 3D transformations
 * - Z-depth sorting of triangles for proper rendering order
 * - Screen buffer management with double buffering
 * - Clamping values to specified ranges
 * - Managing screen dimensions and rotation state
 */

import (
	"math"
	"sort"
)

type Matrix struct {
	cols         int
	rows         int
	angle        float64
	ScreenBuffer [2][][]rune
}

func NewMatrix(cols, rows int) *Matrix {
	m := &Matrix{
		cols:  cols,
		rows:  rows,
		angle: 0,
	}

	m.ScreenBuffer[0] = make([][]rune, rows)
	m.ScreenBuffer[1] = make([][]rune, rows)
	for i := 0; i < rows; i++ {
		m.ScreenBuffer[0][i] = make([]rune, cols+1)
		m.ScreenBuffer[1][i] = make([]rune, cols+1)
	}

	return m
}

func Clamp(value float64, min, max int) int {
	return int(math.Min(float64(max), math.Max(float64(min), value)))
}

func MultiplyMatrices(matrixA, matrixB [][]float64) [][]float64 {
	result := make([][]float64, len(matrixA))
	for i := range result {
		result[i] = make([]float64, len(matrixB[0]))
	}

	for i := 0; i < len(matrixA); i++ {
		for j := 0; j < len(matrixB[0]); j++ {
			for k := 0; k < len(matrixA[0]); k++ {
				result[i][j] += matrixA[i][k] * matrixB[k][j]
			}
		}
	}
	return result
}

type Triangle struct {
	verts [4][]float64
	avgZ  float64
}

func (m *Matrix) SortVerts(verts [][]float64) [][]float64 {
	var triangles []Triangle

	for i := 0; i < len(verts); i += 4 {
		if i+3 < len(verts) {
			triangle := Triangle{
				verts: [4][]float64{verts[i], verts[i+1], verts[i+2], verts[i+3]},
			}

			triangle.avgZ = (verts[i][2] + verts[i+1][2] + verts[i+2][2]) / 3.0
			triangles = append(triangles, triangle)
		}
	}

	sort.Slice(triangles, func(i, j int) bool {
		return triangles[i].avgZ > triangles[j].avgZ
	})

	result := make([][]float64, len(verts))
	idx := 0
	for _, triangle := range triangles {
		for _, vert := range triangle.verts {
			if idx < len(result) {
				result[idx] = vert
				idx++
			}
		}
	}

	return result
}

func (m *Matrix) GetCols() int {
	return m.cols
}

func (m *Matrix) GetRows() int {
	return m.rows
}

func (m *Matrix) GetAngle() float64 {
	return m.angle
}

func (m *Matrix) SetAngle(angle float64) {
	m.angle = angle
}
