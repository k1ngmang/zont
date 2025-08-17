package main

import (
	"log"
)

func main() {
	file := "models/cube.obj"

	matrix := NewMatrix(10, 10)

	renderer := NewRender(matrix)

	verts, err := loadOBJ(file)
	if err != nil {
		log.Fatalf("Error %s: %v", file, err)
	}

	renderer.Render(verts)
}
