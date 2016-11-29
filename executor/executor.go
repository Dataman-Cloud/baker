package executor

import (
	"sync"
	"time"

	_ "golang.org/x/net/context"
)

const timeout = 3 * time.Minute

type Executor struct {
	Pool      *WorkPool
	Works     []*Work // Workpool.
	Collector *Collector
}

type Work struct {
	ID    string
	Tasks []func()
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
		tasks := work.Tasks
		w := func() error {
			defer wg.Done()
			for _, task := range tasks {
				task()
			}
			return nil
		}
		// submit work
		t.Pool.Submit(w)
	}
	wg.Wait()
}
