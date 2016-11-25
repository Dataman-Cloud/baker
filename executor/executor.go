package executor

import (
	"sync"
)

type Executor struct {
	Pool  *WorkPool
	Tasks []Task
}

type Task struct {
	ID   string
	Work func()
}

type TaskStats struct {
	Status chan int32
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

func NewExecutor(pool *WorkPool, tasks []Task) (*Executor, error) {
	return &Executor{
		Pool:  pool,
		Tasks: tasks,
	}, nil
}

func (t *Executor) Execute() {
	wg := sync.WaitGroup{}
	wg.Add(len(t.Tasks))
	for _, task := range t.Tasks {
		work := task.Work
		t.Pool.Submit(func() {
			defer wg.Done()
			work()
		})
	}
	wg.Wait()
}
