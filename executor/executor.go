package executor

import (
	"sync"
)

type Executor struct {
	Pool      *WorkPool
	Tasks     []*Task
	Collector *Collector
}

type Task struct {
	ID   string
	Work func()
}

const (
	StatusStarting = 0
	StatusRunning  = 1
	StatusFailed   = 2
	StatusExpired  = 3
	StatusFinished = 4
)

func NewExecutor(pool *WorkPool, tasks []*Task, collector *Collector) (*Executor, error) {
	return &Executor{
		Pool:      pool,
		Tasks:     tasks,
		Collector: collector,
	}, nil
}

func (t *Executor) Execute() {
	wg := sync.WaitGroup{}
	wg.Add(len(t.Tasks))
	for _, task := range t.Tasks {
		work := task.Work
		task.Work = func() {
			defer wg.Done()
			work()
		}
		// start task collector
		t.Collector.Start()
	}
	wg.Wait()
}
