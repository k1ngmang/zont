package main

import (
	"fmt"
	"log"
	"zontengine/internal/matrix"
	"zontengine/internal/render"
)

func main() {
	file := "models/cube.obj"
	matrix := matrix.NewMatrix(20, 20)

	renderer := render.NewRender(matrix)

	verts, err := render.LoadOBJ(file)
	if err != nil {
		log.Fatalf("Error %s: %v", file, err)
	}

	str := renderer.RenderFrontFace(verts)
	fmt.Println(str)
}
