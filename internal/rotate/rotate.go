package rotate

/**
 * Manages 3D rotation matrices for transforming 3D vertices.
 * Stores and provides access to X, Y, and Z axis rotation matrices.
 *
 * The rotation matrices are used to transform 3D coordinates by:
 * - X-axis rotation: pitch (rotation around horizontal axis)
 * - Y-axis rotation: yaw (rotation around vertical axis)
 * - Z-axis rotation: roll (rotation around depth axis)
 *
 * These matrices are typically updated from the renderer based on
 * the current rotation angle and applied to vertices during processing.
 */

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
