package executor

import (
	"sync"
	"time"

	"github.com/Sirupsen/logrus"
	"golang.org/x/net/context"
)

const timeout = 5 * time.Minute

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
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	wg := sync.WaitGroup{}
	wg.Add(len(t.Works))
	for _, work := range t.Works {
		tasks := work.Tasks
		w := func(ctx context.Context) error {
			defer wg.Done()
			select {
			case <-ctx.Done():
				logrus.Info("Cancel the context")
				return ctx.Err()
			default:
				for _, task := range tasks {
					task()
				}
			}
			return nil
		}
		// submit work
		t.Pool.Submit(ctx, w)
	}
	wg.Wait()
}
