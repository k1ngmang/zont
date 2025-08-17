package main

import (
	"math"
	"sort"
)

type Matrix struct {
	cols         int
	rows         int
	angle        float64
	screenBuffer [2][][]rune
}

func NewMatrix(cols, rows int) *Matrix {
	m := &Matrix{
		cols:  cols,
		rows:  rows,
		angle: 0,
	}

	m.screenBuffer[0] = make([][]rune, rows)
	m.screenBuffer[1] = make([][]rune, rows)
	for i := 0; i < rows; i++ {
		m.screenBuffer[0][i] = make([]rune, cols+1)
		m.screenBuffer[1][i] = make([]rune, cols+1)
	}

	return m
}

func clamp(value float64, min, max int) int {
	return int(math.Min(float64(max), math.Max(float64(min), value)))
}

func multiplyMatrices(matrixA, matrixB [][]float64) [][]float64 {
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

func (m *Matrix) sortVerts(verts [][]float64) [][]float64 {
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
