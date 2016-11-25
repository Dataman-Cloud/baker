package executor

import (
	"sync"
)

type Executor struct {
	Pool  *WorkPool
	Tasks []*Task
}

type Task struct {
	ID     string
	Work   func()
	Status chan int
}

const (
	StatusStarting         = 0
	StatusRunning          = 1
	StatusFailed           = 2
	StatusExpired          = 3
	StatusFinished         = 4
	StatusDockerLoginStart = 5
	StatusDockerLoginOK    = 6
	StatusDockerBuildStart = 7
	StatusDockerBuildOK    = 8
	StatusDockerPushStart  = 9
	StatusDockerPushOK     = 10
)

func NewExecutor(pool *WorkPool, tasks []*Task) (*Executor, error) {
	return &Executor{
		Pool:  pool,
		Tasks: tasks,
	}, nil
}

func (t *Executor) Execute() {
	wg := sync.WaitGroup{}
	wg.Add(len(t.Tasks))
	for _, task := range t.Tasks {
		t.Pool.Submit(task)
		defer wg.Done()
	}
	wg.Wait()
}
