package executor

import (
	"fmt"
	"sync"
)

type Executor struct {
	pool  *WorkPool
	works []func()
}

func NewExecutor(maxWorkers int, works []func()) (*Executor, error) {
	if maxWorkers < 1 {
		return nil, fmt.Errorf("must provide positive maxWorkers; provided %d", maxWorkers)
	}

	var pool *WorkPool
	if len(works) < maxWorkers {
		pool = newWorkPoolWithPending(len(works), 0)
	} else {
		pool = newWorkPoolWithPending(maxWorkers, len(works)-maxWorkers)
	}

	return &Executor{
		pool:  pool,
		works: works,
	}, nil
}

func (t *Executor) Execute() {
	defer t.pool.Stop()

	wg := sync.WaitGroup{}
	wg.Add(len(t.works))
	for _, work := range t.works {
		work := work
		t.pool.Submit(func() {
			defer wg.Done()
			work()
		})
	}
	wg.Wait()
}
