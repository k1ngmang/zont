package rotate

type Rotate struct {
	x [][]float64
	y [][]float64
	z [][]float64
}

func NewRotate() *Rotate {
	return &Rotate{}
}

func (r *Rotate) Update(x, y, z [][]float64) {
	r.x = x
	r.y = y
	r.z = z
}

func (r *Rotate) GetX() [][]float64 {
	return r.x
}

func (r *Rotate) GetY() [][]float64 {
	return r.y
}

func (r *Rotate) GetZ() [][]float64 {
	return r.z
}
