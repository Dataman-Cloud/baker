package executor

import (
	"sync"
)

type Executor struct {
	Pool  *WorkPool
	Works []func()
}

func NewExecutor(pool *WorkPool, works []func()) (*Executor, error) {
	return &Executor{
		Pool:  pool,
		Works: works,
	}, nil
}

func (t *Executor) Execute() {
	defer t.Pool.Stop()

	wg := sync.WaitGroup{}
	wg.Add(len(t.Works))
	for _, work := range t.Works {
		work := work
		t.Pool.Submit(func() {
			defer wg.Done()
			work()
		})
	}
	wg.Wait()
}
