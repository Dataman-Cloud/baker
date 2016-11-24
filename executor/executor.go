package executor

import (
	"sync"

	"github.com/Sirupsen/logrus"
)

type Executor struct {
	Pool  *WorkPool
	Tasks []Task
}

type Task struct {
	Name   string
	Work   func()
	Status int
}

func NewExecutor(pool *WorkPool, tasks []Task) (*Executor, error) {
	return &Executor{
		Pool:  pool,
		Tasks: tasks,
	}, nil
}

func (t *Executor) Execute() {
	defer t.Pool.Stop()

	wg := sync.WaitGroup{}
	wg.Add(len(t.Tasks))
	for _, task := range t.Tasks {
		logrus.Info("WAITING")
		work := task.Work
		t.Pool.Submit(func() {
			defer wg.Done()
			defer func() {
				if r := recover(); r != nil {
					logrus.Info("FAILED")
				} else {
					logrus.Info("FINISHED")
				}
			}()
			work()
		})
	}
	wg.Wait()
}
