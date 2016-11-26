package executor

import (
	"github.com/Sirupsen/logrus"
)

type Collector struct {
	TaskID    string
	TaskStats chan int
	TaskMsg   chan string
	IsDone    chan bool
}

func NewCollector(taskID string, taskStats chan int, taskMsg chan string, isDone chan bool) *Collector {
	return &Collector{
		TaskID:    taskID,
		TaskStats: taskStats,
		TaskMsg:   taskMsg,
		IsDone:    isDone,
	}
}

func (c *Collector) Start() {
	go func() {
		logrus.Info("Collector start")
		for {
			select {
			case ts := <-c.TaskStats:
				logrus.Info("taskID:%s status:%s", c.TaskID, ts)
				if ts == StatusFinished || ts == StatusFailed || ts == StatusExpired {
					c.IsDone <- true
				}
			case tm := <-c.TaskMsg:
				logrus.Info("taskID:%s message:%s", c.TaskID, tm)
			case <-c.IsDone:
				break
			}
		}
	}()
}
