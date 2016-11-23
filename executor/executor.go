package executor

import (
	"errors"
	"sync"

	"github.com/Sirupsen/logrus"
)

type Executor struct {
	Pool  *WorkPool
	Works []func()
}

func NewExecutor(maxWorkers int, works []func()) (*Executor, error) {
	if maxWorkers < 1 {
		logrus.Fatalf("must provide positive maxWorkers; provided %d", maxWorkers)

		return nil, errors.New("must provide positive maxWorkers")
	}
	var pool *WorkPool
	if len(works) < maxWorkers {
		pool = newWorkPoolWithPending(len(works), 0)
	} else {
		pool = newWorkPoolWithPending(maxWorkers, len(works)-maxWorkers)
	}

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
