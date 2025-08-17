package main

func toArray1D(array2D [][]float64) []float64 {
	array1D := make([]float64, len(array2D))
	for i := 0; i < len(array1D); i++ {
		array1D[i] = array2D[i][0]
	}
	return array1D
}

func toArray2D(array1D []float64) [][]float64 {
	array2D := make([][]float64, len(array1D))
	for i := 0; i < len(array2D); i++ {
		array2D[i] = []float64{array1D[i]}
	}
	return array2D
}
