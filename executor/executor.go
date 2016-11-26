package executor

import (
	"sync"
)

type Executor struct {
	Pool      *WorkPool
	Works     []*Work
	Collector *Collector
}

type Work struct {
	ID   string
	Task func()
}

func NewExecutor(pool *WorkPool, works []*Work, collector *Collector) (*Executor, error) {
	return &Executor{
		Pool:      pool,
		Works:     works,
		Collector: collector,
	}, nil
}

func (t *Executor) Execute() {
	// defer t.Pool.Stop() // stop pool in baker server stop.
	wg := sync.WaitGroup{}
	wg.Add(len(t.Works))
	for _, work := range t.Works {
		task := work.Task
		work.Task = func() {
			defer wg.Done()
			task()
		}
		// start task collector
		t.Collector.Start()
		t.Collector.TaskStatus <- StatusStarting
		// submit work
		t.Pool.Submit(work)
	}
	wg.Wait()
}
