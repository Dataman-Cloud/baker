package executor

import (
	"fmt"
	"sync"
	"sync/atomic"
)

type WorkPool struct {
	workQueue chan func() error
	stopping  chan struct{}
	stopped   int32

	mutex       sync.Mutex
	maxWorkers  int
	numWorkers  int
	idleWorkers int
}

func NewWorkPool(maxWorkers int) (*WorkPool, error) {
	if maxWorkers < 1 {
		return nil, fmt.Errorf("must provide positive maxWorkers; provided %d", maxWorkers)
	}

	return newWorkPoolWithPending(maxWorkers, 0), nil
}

func newWorkPoolWithPending(maxWorkers, pending int) *WorkPool {
	return &WorkPool{
		workQueue:  make(chan func() error, maxWorkers+pending),
		stopping:   make(chan struct{}),
		maxWorkers: maxWorkers,
	}
}

func (w *WorkPool) Submit(work func() error) {
	if atomic.LoadInt32(&w.stopped) == 1 {
		return
	}

	select {
	case w.workQueue <- work:
		if atomic.LoadInt32(&w.stopped) == 1 {
			w.drain()
		} else {
			w.addWorker()
		}
	case <-w.stopping:
	}
}

func (w *WorkPool) Stop() {
	if atomic.CompareAndSwapInt32(&w.stopped, 0, 1) {
		close(w.stopping)
		w.drain()
	}
}

func (w *WorkPool) addWorker() bool {
	w.mutex.Lock()
	defer w.mutex.Unlock()

	if w.idleWorkers > 0 || w.numWorkers == w.maxWorkers {
		return false
	}

	w.numWorkers++
	go worker(w)
	return true
}

func (w *WorkPool) workerStopping(force bool) bool {
	w.mutex.Lock()
	if !force {
		if len(w.workQueue) <= w.numWorkers {
			w.mutex.Unlock()
			return false
		}
	}

	w.numWorkers--
	w.mutex.Unlock()

	return true
}

func (w *WorkPool) drain() {
	for {
		select {
		case <-w.workQueue:
		default:
			return
		}
	}
}

func worker(w *WorkPool) {
	for {
		if atomic.LoadInt32(&w.stopped) == 1 {
			w.workerStopping(true)
			return
		}
		select {
		case <-w.stopping:
			w.workerStopping(true)
			return
		case work := <-w.workQueue:
			w.mutex.Lock()
			w.idleWorkers--
			w.mutex.Unlock()

		NOWORK:
			for {
				work()
				select {
				case work = <-w.workQueue:
				case <-w.stopping:
					break NOWORK
				default:
					break NOWORK
				}
			}
			w.mutex.Lock()
			w.idleWorkers++
			w.mutex.Unlock()
		}
	}
}
