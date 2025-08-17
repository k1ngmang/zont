package main

type Rotate struct {
	x [][]float64
	y [][]float64
	z [][]float64
}

func NewRotate() *Rotate {
	return &Rotate{}
}

func (r *Rotate) update(x, y, z [][]float64) {
	r.x = x
	r.y = y
	r.z = z
}

func (r *Rotate) getX() [][]float64 {
	return r.x
}

func (r *Rotate) getY() [][]float64 {
	return r.y
}

func (r *Rotate) getZ() [][]float64 {
	return r.z
}
