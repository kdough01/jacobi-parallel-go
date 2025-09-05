package jacobi

import (
	"sync"
	"fmt"
)

type SharedContext struct {
	mutex       *sync.Mutex
	cond        *sync.Cond
	wgContext   *sync.WaitGroup
	counter     int
	threadCount int
	xNew		[]float64
}

func Slice(ASlice [][]float64, B []float64, start int, end int, x []float64, ctx *SharedContext) []float64 {
	n := end - start
	xNew := make([]float64, n)

	for i:=start;i<end;i++ {
		sigma := 0.0
		for j:=0;j<len(ASlice[i]);j++ {
			if j != i {
				sigma = sigma + ASlice[i][j]*x[j]
			}
		}
		xNew[i - start] = (B[i] - sigma) / ASlice[i][i]
	}

	return xNew
}

func Worker(goID int, ctx *SharedContext, ASlice [][]float64, B []float64, x []float64, start int, end int) {
	xSlice := Slice(ASlice, B, start, end, x, ctx)

	ctx.mutex.Lock()
	for i:=start;i<end;i++ {
		ctx.xNew[i] = xSlice[i - start]
	}
	ctx.counter++
	if ctx.counter == ctx.threadCount {
		ctx.cond.Broadcast()
	} else {
		for ctx.counter != ctx.threadCount {
			ctx.cond.Wait()
		}
	}
	ctx.mutex.Unlock()

	ctx.wgContext.Done()
}

func BSP(numThreads int, maxIter int, A [][]float64, B []float64, desiredConvergence float64) {

	var wg sync.WaitGroup
	threadCount := numThreads
	n := len(B)
	x := make([]float64, n)
	
	var mutex sync.Mutex
	condVar := sync.NewCond(&mutex)
	context := SharedContext{wgContext: &wg, threadCount: threadCount, cond: condVar, mutex: &mutex, xNew: make([]float64, len(B))}

	for iter:=0;iter<maxIter;iter++ {
		rowsPerThread := len(A) / numThreads
		remainder := len(A) % numThreads
		start := 0
		context.counter = 0
		for r:=0;r<numThreads;r++ {
			end := start + rowsPerThread
			if r < remainder {
				end ++
			}

			wg.Add(1)
			go Worker(r, &context, A, B, x, start, end)
			start = end
		}
		wg.Wait()
		actualConvergence := 0.0
		actualConvergence, x, context.xNew = ComputeConvergence(n, x, context.xNew)

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