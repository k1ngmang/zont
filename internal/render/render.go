package render

import (
	"bufio"
	"math"
	"os"
	"strconv"
	"strings"
	"time"
	"zontengine/internal/convert"
	"zontengine/internal/matrix"
	"zontengine/internal/rotate"
	"zontengine/internal/screen"
)

type Render struct {
	matrix *matrix.Matrix
	screen *screen.Screen
	rotate *rotate.Rotate
}

func NewRender(matrix *matrix.Matrix) *Render {
	return &Render{
		matrix: matrix,
		screen: screen.NewScreen(matrix),
		rotate: rotate.NewRotate(),
	}
}

func (r *Render) Render(verts [][]float64) {
	go r.renderThread()

	for {
		r.updateRotation()
		visibleVerts := r.processVertices(verts)
		vertsToRender := r.matrix.SortVerts(visibleVerts)

		r.screen.InitScreen(r.matrix.ScreenBuffer[0])

		for i := 0; i < len(vertsToRender); i += 4 {
			if i+3 < len(vertsToRender) {
				r.fillTriangle(vertsToRender[i], vertsToRender[i+1], vertsToRender[i+2], vertsToRender[i+3])
			}
		}

		for i := 0; i < len(r.matrix.ScreenBuffer[0]); i++ {
			copy(r.matrix.ScreenBuffer[1][i], r.matrix.ScreenBuffer[0][i])
		}
	}
}

func (r *Render) RenderFrontFace(verts [][]float64) string {
	tempBuffer := make([][]rune, len(r.matrix.ScreenBuffer[0]))
	for i := range tempBuffer {
		tempBuffer[i] = make([]rune, len(r.matrix.ScreenBuffer[0][0]))
		for j := range tempBuffer[i] {
			tempBuffer[i][j] = ' '
		}
	}

	originalAngle := r.matrix.GetAngle()
	r.matrix.SetAngle(0)
	r.updateRotation()

	visibleVerts := r.processVertices(verts)
	vertsToRender := r.matrix.SortVerts(visibleVerts)

	for i := 0; i < len(vertsToRender); i += 4 {
		if i+3 < len(vertsToRender) {
			r.fillTriangleToBuffer(tempBuffer, vertsToRender[i], vertsToRender[i+1], vertsToRender[i+2], vertsToRender[i+3])
		}
	}

	r.matrix.SetAngle(originalAngle)

	var result strings.Builder
	for i := 0; i < len(tempBuffer); i++ {
		for j := 0; j < len(tempBuffer[i]); j++ {
			result.WriteRune(tempBuffer[i][j])
		}
		if i < len(tempBuffer)-1 {
			result.WriteString("\n")
		}
	}

	return result.String()
}

func (r *Render) fillTriangleToBuffer(buffer [][]rune, vert1, vert2, vert3, normal []float64) {
	tempScreen := make([][]rune, len(buffer))
	for i := range tempScreen {
		tempScreen[i] = make([]rune, len(buffer[0]))
		for j := range tempScreen[i] {
			tempScreen[i][j] = ' '
		}
	}

	projection := [][]float64{
		{1, 0, 0},
		{0, 1, 0},
	}

	shadingChars := []rune{'.', ',', '-', '~', ':', ';', '=', '!', '*', '#', '$', '@'}

	lightDirection := []float64{0, 0, -1}
	magnitude := math.Sqrt(lightDirection[0]*lightDirection[0] + lightDirection[1]*lightDirection[1] + lightDirection[2]*lightDirection[2])
	lightDirection[0] /= magnitude
	lightDirection[1] /= magnitude
	lightDirection[2] /= magnitude

	dot := normal[0]*lightDirection[0] + normal[1]*lightDirection[1] + normal[2]*lightDirection[2]
	shadingChar := shadingChars[matrix.Clamp(dot*12, 0, len(shadingChars)-1)]

	vert1 = convert.ToArray1D(matrix.MultiplyMatrices(projection, convert.ToArray2D(vert1)))
	vert2 = convert.ToArray1D(matrix.MultiplyMatrices(projection, convert.ToArray2D(vert2)))
	vert3 = convert.ToArray1D(matrix.MultiplyMatrices(projection, convert.ToArray2D(vert3)))

	r.drawLineToBuffer(tempScreen, vert1[0], vert1[1], vert2[0], vert2[1], shadingChar)
	r.drawLineToBuffer(tempScreen, vert2[0], vert2[1], vert3[0], vert3[1], shadingChar)
	r.drawLineToBuffer(tempScreen, vert3[0], vert3[1], vert1[0], vert1[1], shadingChar)

	r.fillTriangleArea(tempScreen, shadingChar)

	for i := 0; i < len(buffer); i++ {
		for j := 0; j < len(buffer[0]); j++ {
			if tempScreen[i][j] == shadingChar {
				buffer[i][j] = tempScreen[i][j]
			}
		}
	}
}

func (r *Render) drawLineToBuffer(buffer [][]rune, x1, y1, x2, y2 float64, ch rune) {
	x1 = float64(r.matrix.GetCols())/2.0 + x1/2.0*float64(r.matrix.GetCols())
	y1 = float64(r.matrix.GetRows())/2.0 + y1/-2.0*float64(r.matrix.GetRows())
	x2 = float64(r.matrix.GetCols())/2.0 + x2/2.0*float64(r.matrix.GetCols())
	y2 = float64(r.matrix.GetRows())/2.0 + y2/-2.0*float64(r.matrix.GetRows())

	d := 0
	dx := int(math.Abs(x2 - x1))
	dy := int(math.Abs(y2 - y1))

	dx2 := 2 * dx
	dy2 := 2 * dy

	ix := 1
	if x1 > x2 {
		ix = -1
	}
	iy := 1
	if y1 > y2 {
		iy = -1
	}

	x := int(x1)
	y := int(y1)

	if dx >= dy {
		for {
			if y >= 0 && y < len(buffer) && x >= 0 && x < len(buffer[0]) && buffer[y][x] == ' ' {
				buffer[y][x] = ch
			}
			if x == int(x2) {
				break
			}
			x += ix
			d += dy2
			if d > dx {
				y += iy
				d -= dx2
			}
		}
	} else {
		for {
			if y >= 0 && y < len(buffer) && x >= 0 && x < len(buffer[0]) && buffer[y][x] == ' ' {
				buffer[y][x] = ch
			}
			if y == int(y2) {
				break
			}
			y += iy
			d += dx2
			if d > dy {
				x += ix
				d -= dy2
			}
		}
	}
}

func (r *Render) renderThread() {
	fps := 60
	for {
		start := time.Now()
		r.screen.DrawScreen()
		r.matrix.SetAngle(r.matrix.GetAngle() + 0.03*(60.0/float64(fps)))

		elapsed := time.Since(start)
		sleepTime := time.Duration(1000/fps)*time.Millisecond - elapsed
		if sleepTime > 0 {
			time.Sleep(sleepTime)
		}
	}
}

func (r *Render) updateRotation() {
	angle := r.matrix.GetAngle()

	xMatrix := [][]float64{
		{1, 0, 0},
		{0, 1, 0},
		{0, 0, 1},
	}

	yMatrix := [][]float64{
		{math.Cos(angle), 0, math.Sin(angle)},
		{0, 1, 0},
		{-math.Sin(angle), 0, math.Cos(angle)},
	}

	zMatrix := [][]float64{
		{1, 0, 0},
		{0, 1, 0},
		{0, 0, 1},
	}

	r.rotate.Update(xMatrix, yMatrix, zMatrix)
}

func (r *Render) processVertices(verts [][]float64) [][]float64 {
	var visibleVerts [][]float64

	for i := 0; i < len(verts); i += 3 {
		if i+2 >= len(verts) {
			break
		}

		vert1 := convert.ToArray1D(matrix.MultiplyMatrices(r.rotate.GetX(), matrix.MultiplyMatrices(r.rotate.GetY(), matrix.MultiplyMatrices(r.rotate.GetZ(), convert.ToArray2D(verts[i])))))
		vert2 := convert.ToArray1D(matrix.MultiplyMatrices(r.rotate.GetX(), matrix.MultiplyMatrices(r.rotate.GetY(), matrix.MultiplyMatrices(r.rotate.GetZ(), convert.ToArray2D(verts[i+1])))))
		vert3 := convert.ToArray1D(matrix.MultiplyMatrices(r.rotate.GetX(), matrix.MultiplyMatrices(r.rotate.GetY(), matrix.MultiplyMatrices(r.rotate.GetZ(), convert.ToArray2D(verts[i+2])))))

		normal := []float64{
			((vert2[1] - vert1[1]) * (vert3[2] - vert1[2])) - ((vert2[2] - vert1[2]) * (vert3[1] - vert1[1])),
			((vert2[2] - vert1[2]) * (vert3[0] - vert1[0])) - ((vert2[0] - vert1[0]) * (vert3[2] - vert1[2])),
			((vert2[0] - vert1[0]) * (vert3[1] - vert1[1])) - ((vert2[1] - vert1[1]) * (vert3[0] - vert1[0])),
		}

		magnitude := math.Sqrt(normal[0]*normal[0] + normal[1]*normal[1] + normal[2]*normal[2])
		if magnitude > 0 {
			normal[0] /= magnitude
			normal[1] /= magnitude
			normal[2] /= magnitude
		}

		if normal[0]*vert1[0]+normal[1]*vert1[1]+normal[2]*(vert1[2]-10) > 1 {
			visibleVerts = append(visibleVerts, vert1, vert2, vert3, normal)
		}
	}

	return visibleVerts
}

func (r *Render) fillTriangle(vert1, vert2, vert3, normal []float64) {
	tempScreen := make([][]rune, len(r.matrix.ScreenBuffer[0]))
	for i := range tempScreen {
		tempScreen[i] = make([]rune, len(r.matrix.ScreenBuffer[0][0]))
	}
	r.screen.InitScreen(tempScreen)

	projection := [][]float64{
		{1, 0, 0},
		{0, 1, 0},
	}

	shadingChars := []rune{'.', ',', '-', '~', ':', ';', '=', '!', '*', '#', '$', '@'}

	lightDirection := []float64{0, 0, -1}
	magnitude := math.Sqrt(lightDirection[0]*lightDirection[0] + lightDirection[1]*lightDirection[1] + lightDirection[2]*lightDirection[2])
	lightDirection[0] /= magnitude
	lightDirection[1] /= magnitude
	lightDirection[2] /= magnitude

	dot := normal[0]*lightDirection[0] + normal[1]*lightDirection[1] + normal[2]*lightDirection[2]
	shadingChar := shadingChars[matrix.Clamp(dot*12, 0, len(shadingChars)-1)]

	vert1 = convert.ToArray1D(matrix.MultiplyMatrices(projection, convert.ToArray2D(vert1)))
	vert2 = convert.ToArray1D(matrix.MultiplyMatrices(projection, convert.ToArray2D(vert2)))
	vert3 = convert.ToArray1D(matrix.MultiplyMatrices(projection, convert.ToArray2D(vert3)))

	r.drawLine(tempScreen, vert1[0], vert1[1], vert2[0], vert2[1], shadingChar)
	r.drawLine(tempScreen, vert2[0], vert2[1], vert3[0], vert3[1], shadingChar)
	r.drawLine(tempScreen, vert3[0], vert3[1], vert1[0], vert1[1], shadingChar)

	r.fillTriangleArea(tempScreen, shadingChar)

	for i := 0; i < len(r.matrix.ScreenBuffer[0]); i++ {
		for j := 0; j < len(r.matrix.ScreenBuffer[0][0]); j++ {
			if tempScreen[i][j] == shadingChar {
				r.matrix.ScreenBuffer[0][i][j] = tempScreen[i][j]
			}
		}
	}
}

func (r *Render) drawLine(screen [][]rune, x1, y1, x2, y2 float64, ch rune) {
	x1 = float64(r.matrix.GetCols())/2.0 + x1/2.0*float64(r.matrix.GetCols())
	y1 = float64(r.matrix.GetRows())/2.0 + y1/-2.0*float64(r.matrix.GetRows())
	x2 = float64(r.matrix.GetCols())/2.0 + x2/2.0*float64(r.matrix.GetCols())
	y2 = float64(r.matrix.GetRows())/2.0 + y2/-2.0*float64(r.matrix.GetRows())

	d := 0
	dx := int(math.Abs(x2 - x1))
	dy := int(math.Abs(y2 - y1))

	dx2 := 2 * dx
	dy2 := 2 * dy

	ix := 1
	if x1 > x2 {
		ix = -1
	}
	iy := 1
	if y1 > y2 {
		iy = -1
	}

	x := int(x1)
	y := int(y1)

	if dx >= dy {
		for {
			if y >= 0 && y < len(screen) && x >= 0 && x < len(screen[0]) && screen[y][x] == ' ' {
				screen[y][x] = ch
			}
			if x == int(x2) {
				break
			}
			x += ix
			d += dy2
			if d > dx {
				y += iy
				d -= dx2
			}
		}
	} else {
		for {
			if y >= 0 && y < len(screen) && x >= 0 && x < len(screen[0]) && screen[y][x] == ' ' {
				screen[y][x] = ch
			}
			if y == int(y2) {
				break
			}
			y += iy
			d += dx2
			if d > dy {
				x += ix
				d -= dy2
			}
		}
	}
}

func (r *Render) fillTriangleArea(screen [][]rune, shadingChar rune) {
	for row := 0; row < len(screen); row++ {
		rowStr := string(screen[row])
		shadingStr := string(shadingChar)
		firstIndex := strings.Index(rowStr, shadingStr)
		lastIndex := strings.LastIndex(rowStr, shadingStr)

		if firstIndex != -1 && lastIndex != -1 && firstIndex < lastIndex {
			for i := firstIndex; i < lastIndex; i++ {
				screen[row][i] = shadingChar
			}
		}
	}
}

func LoadOBJ(filename string) ([][]float64, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var verts [][]float64
	var faces []int

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		if len(line) == 0 {
			continue
		}

		if line[0] == 'v' {
			parts := strings.Fields(line)
			if len(parts) >= 4 {
				x, _ := strconv.ParseFloat(parts[1], 64)
				y, _ := strconv.ParseFloat(parts[2], 64)
				z, _ := strconv.ParseFloat(parts[3], 64)
				verts = append(verts, []float64{x, y, z})
			}
		} else if line[0] == 'f' {
			parts := strings.Fields(line)
			if len(parts) >= 4 {
				f1, _ := strconv.Atoi(parts[1])
				f2, _ := strconv.Atoi(parts[2])
				f3, _ := strconv.Atoi(parts[3])
				faces = append(faces, f1-1, f2-1, f3-1) // OBJ индексы начинаются с 1
			}
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	result := make([][]float64, len(faces))
	for i := 0; i < len(faces); i++ {
		if faces[i] < len(verts) {
			result[i] = verts[faces[i]]
		}
	}

	return result, nil
}
