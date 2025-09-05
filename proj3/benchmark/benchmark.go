package main

import (
	"math/rand"
	"proj3-redesigned/jacobi"
	"time"
	"math"
	"flag"
	"fmt"
)

func randomMatrix(rows int, cols int) [][]float64 {
	rand.Seed(time.Now().Unix())
	mat := make([][]float64, rows)

	for col := range mat {
		mat[col] = make([]float64, cols)
	}

	for r:=0;r<rows;r++ {
		sum := 0.0
		for c:=0;c<cols;c++ {
			mat[r][c] = float64(rand.Intn(1000))
			sum += math.Abs(mat[r][c])
		}
		mat[r][r] = sum + float64(rand.Intn(10))
	}

	return mat
}

func randomArray(size int) []float64 {
	rand.Seed(time.Now().Unix())
	arr := make([]float64, size)

	for s:=0;s<size;s++ {
		arr[s] = float64(rand.Intn(1000))
	}

	return arr
}

func main() {

	size := flag.Int("size", 100, "Size of the matrix (NxN)")
	rowsPerTask := flag.Int("rowsPerTask", 2, "Number of rows per task for work stealing")
	numThreads := flag.Int("numThreads", 2, "Number of threads to use")
	desiredConvergence := flag.Float64("convergence", 0.000001, "Desired convergence threshold")
	maxIter := flag.Int("maxIter", 10000, "Maximum number of iterations")
	mode := flag.String("mode", "sequential", "Execution mode: sequential, bsp, or worksteal")

	flag.Parse()

	A := randomMatrix(*size, *size)
	B := randomArray(*size)

	fmt.Printf("Running with size=%d, rowsPerTask=%d, numThreads=%d, convergence=%.8f, maxIter=%d\n",
		*size, *rowsPerTask, *numThreads, *desiredConvergence, *maxIter)

	start := time.Now()

	switch *mode {
	case "sequential":
		jacobi.Sequential(A, B, *desiredConvergence, *maxIter)
	case "bsp":
		jacobi.BSP(*numThreads, *maxIter, A, B, *desiredConvergence)
	case "worksteal":
		jacobi.WorkSteal(*numThreads, *rowsPerTask, *maxIter, A, B, *desiredConvergence)
	default:
		jacobi.Sequential(A, B, *desiredConvergence, *maxIter)
	}

	elapsed := time.Since(start)
	fmt.Printf("%.6f\n", elapsed.Seconds())
}