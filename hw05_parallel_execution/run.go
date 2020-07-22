package hw05_parallel_execution //nolint:golint,stylecheck

import (
	"errors"
	"sync"
)

// ErrErrorsLimitExceeded ...
var ErrErrorsLimitExceeded = errors.New("errors limit exceeded")

// ErrWrongArgument ...
var ErrWrongArgument = errors.New("wrong argument")

// Task ...
type Task func() error

type syncCounter struct {
	sync.RWMutex
	value int
}

func (e *syncCounter) Inc() {
	e.Lock()
	e.value++
	e.Unlock()
}

func (e *syncCounter) Get() int {
	e.RLock()
	c := e.value
	e.RUnlock()

	return c
}

func runTask(task Task, workers chan struct{}, counter *syncCounter, wg *sync.WaitGroup) {
	defer wg.Done()
	err := task()
	if err != nil {
		counter.Inc()
	}
	<-workers
}

// Run starts tasks in N goroutines and stops its work when receiving M errors from tasks.
func Run(tasks []Task, n int, m int) error {
	if n <= 0 || m < 0 {
		return ErrWrongArgument
	}
	workers := make(chan struct{}, n)
	defer close(workers)
	wg := sync.WaitGroup{}
	defer wg.Wait()
	errorsCounter := syncCounter{}

	for _, task := range tasks {
		if errorsCounter.Get() >= m {
			return ErrErrorsLimitExceeded
		}
		workers <- struct{}{}
		wg.Add(1)
		go runTask(task, workers, &errorsCounter, &wg)
	}

	return nil
}
