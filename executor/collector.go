package executor

import (
	"errors"
	"net/http"

	"github.com/Sirupsen/logrus"
	"github.com/gin-gonic/gin"
	"github.com/manucorporat/sse"
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

func (c *Collector) Stream(ctx *gin.Context, ssEvent *sse.Event) {
	w := ctx.Writer
	clientClose := w.CloseNotify()
	var clientClosed bool = false
	go func() {
		logrus.Info("Collector start")
		for {
			select {
			case <-clientClose:
				clientClosed = true
				logrus.Infof("Client closed")
			case ts := <-c.TaskStatus:
				TaskStatus := TaskStatusEnum[ts]
				logrus.Infof("taskID:%s status:%s", c.TaskID, TaskStatus)
				if ts == StatusFinished || ts == StatusFailed || ts == StatusExpired {
					c.IsDone <- true
				}
				if !clientClosed {
					ssEvent.Data = TaskStatus
					ssEvent.Render(w)
					w.Flush()
				}
			case tm := <-c.TaskMsg:
				logrus.Infof("taskID:%s message:%s", c.TaskID, tm)
				ctx.AbortWithError(http.StatusBadRequest, errors.New(tm))
				c.IsDone <- true
			case <-c.IsDone:
				break
			}
		}
	}()
}
