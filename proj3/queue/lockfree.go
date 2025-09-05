package queue

import (
	"sync/atomic"
)

type Task any

type CircularArray struct {
	data  []Task
	mask  int32
}

func NewCircularArray(size int32) *CircularArray {
	if size&(size-1) != 0 {
		panic("size must be power of 2")
	}
	return &CircularArray{data: make([]Task, size), mask: size - 1}
}

func (a *CircularArray) Get(i int32) Task {
	return a.data[i&a.mask]
}

func (a *CircularArray) Set(i int32, t Task) {
	a.data[i&a.mask] = t
}

type WorkStealingDeque struct {
	top    atomic.Int32
	bottom int32
	tasks  *CircularArray
}

func NewWorkStealingDeque(size int32) *WorkStealingDeque {
	return &WorkStealingDeque{tasks: NewCircularArray(size)}
}

func (q *WorkStealingDeque) PushBottom(t Task) {
	b := q.bottom
	q.tasks.Set(b, t)
	q.bottom = b + 1
}

func (q *WorkStealingDeque) PopBottom() (Task, bool) {
	b := q.bottom - 1
	q.bottom = b
	t := q.top.Load()
	size := b - t
	if size < 0 {
		q.bottom = t
		var zero Task
		return zero, false
	}
	task := q.tasks.Get(b)
	if size > 0 {
		return task, true
	}
	if q.top.CompareAndSwap(t, t+1) {
		return task, true
	}
	q.bottom = t + 1
	var zero Task
	return zero, false
}

func (q *WorkStealingDeque) StealTop() (Task, bool) {
	for {
		t := q.top.Load()
		b := q.bottom
		size := b - t
		if size <= 0 {
			var zero Task
			return zero, false
		}
		task := q.tasks.Get(t)
		if q.top.CompareAndSwap(t, t+1) {
			return task, true
		}
	}
}