package executor

import (
	"github.com/Sirupsen/logrus"
)

type Collector struct {
	TaskID     string
	TaskStatus chan int
	TaskMsg    chan string
	IsDone     chan bool
}

func NewCollector(taskID string, taskStatus chan int, taskMsg chan string, isDone chan bool) *Collector {
	return &Collector{
		TaskID:     taskID,
		TaskStatus: taskStatus,
		TaskMsg:    taskMsg,
		IsDone:     isDone,
	}
}

func (c *Collector) Start() {
	go func() {
		logrus.Info("Collector start")
		for {
			select {
			case ts := <-c.TaskStatus:
				logrus.Infof("taskID:%s status:%s", c.TaskID, TaskStatusEnum[ts])
				if ts == StatusFinished || ts == StatusFailed || ts == StatusExpired {
					c.IsDone <- true
				}
			case tm := <-c.TaskMsg:
				logrus.Infof("taskID:%s message:%s", c.TaskID, tm)
				c.IsDone <- true
			case <-c.IsDone:
				break
			}
		}
	}()
}

func (c *Collector) stream() {

}
