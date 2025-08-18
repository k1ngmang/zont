package main

import (
	"log"
	"zontengine/internal/matrix"
	"zontengine/internal/render"
)

func main() {
	file := "models/cube.obj"

	matrix := matrix.NewMatrix(10, 10)

	renderer := render.NewRender(matrix)

	verts, err := render.LoadOBJ(file)
	if err != nil {
		log.Fatalf("Error %s: %v", file, err)
	}

	renderer.Render(verts)
}
