<div align="center">
  <img src="https://github.com/k1ngmang/zont/blob/main/branding/icon.png" width="200">

<h2>Zont</h2>
Simple and fast engine for rendering 3D graphics in the terminal
</div>

### How to use?
First of all, you need to install the go compiler, then in the main file select the necessary settings (dimensions, object for rendering, etc.)
```go
package main

import (
	"log"
	"zontengine/internal/matrix"
	"zontengine/internal/render"
)

func main() {
	file := "models/cube.obj"

	matrix := matrix.NewMatrix(30, 30)

	renderer := render.NewRender(matrix)

	verts, err := render.LoadOBJ(file)
	if err != nil {
		log.Fatalf("Error %s: %v", file, err)
	}

	renderer.Render(verts)
}
```

And using 
```
make build
make run
```
you can run the project

#### Important
For now the project is in the development stage, so from the api you will not be able to conveniently specify the rotation matrix and generally work with the code, but all this will be finalized
