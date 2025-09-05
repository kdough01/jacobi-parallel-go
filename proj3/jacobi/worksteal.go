package jacobi

import (
	"proj3-redesigned/queue"
	"sync"
	"fmt"
)

type Task struct {
	Start int
	End   int
	A     [][]float64
	B     []float64
	X     []float64
	Ctx   *SharedContext
}

func Stealer(id int, deques []*queue.WorkStealingDeque, wg *sync.WaitGroup) {
	defer wg.Done()
	mydq := deques[id]
	numThreads := len(deques)

	for {
		var taskInterface interface{}
		var ok bool

		taskInterface, ok = mydq.PopBottom()
		if !ok {
			stolen := false
			for i := 0; i < numThreads; i++ {
				if i == id {
					continue
				}
				taskInterface, ok = deques[i].StealTop()
				if ok {
					stolen = true
					break
				}
			}
			if !stolen {
				return
			}
		}

		task := taskInterface.(*Task)

		xSlice := Slice(task.A, task.B, task.Start, task.End, task.X, task.Ctx)

		task.Ctx.mutex.Lock()
		for i := task.Start; i < task.End; i++ {
			task.Ctx.xNew[i] = xSlice[i-task.Start]
		}
		task.Ctx.mutex.Unlock()
	}
}

func WorkSteal(numThreads int, rowsPerTask int, maxIter int, A [][]float64, B []float64, desiredConvergence float64) {
	n := len(B)
	x := make([]float64, n)
	xNew := make([]float64, n)

	var mutex sync.Mutex
	var wg sync.WaitGroup
	context := SharedContext{threadCount: numThreads, mutex: &mutex, xNew: xNew}

	for iter := 0; iter < maxIter; iter++ {
		deques := make([]*queue.WorkStealingDeque, numThreads)
		for i := 0; i < numThreads; i++ {
			deques[i] = queue.NewWorkStealingDeque(1024)
		}

		threadId := 0
		for start := 0; start < n; start += rowsPerTask {
			end := start + rowsPerTask
			if end > n {
				end = n
			}
			task := Task{Start: start, End: end, A: A, B: B, X: x, Ctx: &context}
			deques[threadId].PushBottom(&task)
			threadId = (threadId + 1) % numThreads
		}

		wg.Add(numThreads)
		for i := 0; i < numThreads; i++ {
			go Stealer(i, deques, &wg)
		}
		wg.Wait()

		actualConvergence := 0.0
		actualConvergence, x, context.xNew = ComputeConvergence(n, x, context.xNew)

		if actualConvergence <= desiredConvergence {
			fmt.Printf("Converged after %d iterations\n", iter+1)
			break
		}
	}

	fmt.Println("Final solution vector:")
	for i := 0; i < n; i++ {
		fmt.Printf("x[%d] = %.6f\n", i, x[i])
	}
}