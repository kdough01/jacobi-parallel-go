package jacobi

import (
	"fmt"
	"math"
)

func ComputeConvergence(n int, x []float64, xNew []float64) (float64, []float64, []float64) {
	actualConvergence := 0.0
	for i:=0;i<n;i++ {
		actualConvergence += math.Abs(xNew[i] - x[i])
		x[i] = xNew[i]
	}
	return actualConvergence, x, xNew
}

func Sequential(A [][]float64, B []float64, desiredConvergence float64, maxIter int) {
	n := len(B)
	x := make([]float64, n)
	xNew := make([]float64, n)

	for iter:=0;iter<maxIter;iter++ {
		for i:=0;i<n;i++ {
			sigma := 0.0
			for j:=0;j<n;j++ {
				if j != i {
					sigma = sigma + A[i][j]*x[j]
				}
			}
			xNew[i] = (B[i] - sigma) / A[i][i]
		}

		actualConvergence := 0.0
		actualConvergence, x, xNew = ComputeConvergence(n, x, xNew)

		if actualConvergence <= desiredConvergence {
			fmt.Printf("Converged after %d iterations\n", iter + 1)
			break
		}
	}
	fmt.Println("Final solution vector:")
	for i := 0; i < n; i++ {
		fmt.Printf("x[%d] = %.6f\n", i, x[i])
	}

}
