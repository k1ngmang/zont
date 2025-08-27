package rotate

/*
!WARNING: not implemented
*/

import (
	"math"
)

type EulerRotate struct {
	x [][]float64
	y [][]float64
	z [][]float64
}

func NewEulerRotate() *EulerRotate {
	r := &EulerRotate{}
	return r
}

func (r *EulerRotate) Update(x, y, z [][]float64) {
	r.x = x
	r.y = y
	r.z = z
}

// Z => Y => X
func (r *EulerRotate) SetEulerAngles(yaw, pitch, roll float64) {
	yawRad := yaw * math.Pi / 180.0
	pitchRad := pitch * math.Pi / 180.0
	rollRad := roll * math.Pi / 180.0

	cosZ := math.Cos(yawRad)
	sinZ := math.Sin(yawRad)
	zMatrix := [][]float64{
		{cosZ, -sinZ, 0},
		{sinZ, cosZ, 0},
		{0, 0, 1},
	}

	cosY := math.Cos(pitchRad)
	sinY := math.Sin(pitchRad)
	yMatrix := [][]float64{
		{cosY, 0, sinY},
		{0, 1, 0},
		{-sinY, 0, cosY},
	}

	cosX := math.Cos(rollRad)
	sinX := math.Sin(rollRad)
	xMatrix := [][]float64{
		{1, 0, 0},
		{0, cosX, -sinX},
		{0, sinX, cosX},
	}

	r.Update(xMatrix, yMatrix, zMatrix)
}

func (r *EulerRotate) SetAxisAngles(aX, aY, aZ, angleDeg float64) {
	angle := angleDeg * math.Pi / 180.0
	c := math.Cos(angle)
	s := math.Sin(angle)
	t := 1 - c

	len := math.Sqrt(aX*aX + aY*aY + aZ*aZ)

	if len > 0 {
		aX /= len
		aY /= len
		aZ /= len
	}

	matrix := [][]float64{
		{
			t*aX*aX + c,
			t*aX*aY + s*aZ,
			t*aX*aZ + s*aY,
		},
		{
			t*aX*aY + s*aZ,
			t*aY*aY + c,
			t*aY*aZ + s*aX,
		},
		{
			t*aX*aZ + s*aY,
			t*aY*aZ + s*aX,
			t*aZ*aZ + c,
		},
	}
	r.Update(matrix, matrix, matrix)
}
