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
}

func NewCollector(taskID string, taskStatus chan int, taskMsg chan string) *Collector {
	return &Collector{
		TaskID:     taskID,
		TaskStatus: taskStatus,
		TaskMsg:    taskMsg,
	}
}

func (c *Collector) Stream(ctx *gin.Context, ssEvent *sse.Event) {
	w := ctx.Writer
	clientClose := w.CloseNotify()
	go func() {
		logrus.Info("Collector start")
		for {
			select {
			case <-clientClose:
				logrus.Infof("Close Nodify")
				return
			case ts := <-c.TaskStatus:
				TaskStatus := TaskStatusEnum[ts]
				logrus.Infof("taskID:%s status:%s", c.TaskID, TaskStatus)
				ssEvent.Data = "taskID:" + c.TaskID + " " + "status:" + TaskStatus
				ssEvent.Render(w)
				w.Flush()
				if ts == StatusFinished || ts == StatusFailed || ts == StatusExpired {
					ssEvent.Data = "CLOSE"
					ssEvent.Render(w)
					w.Flush()
				}
			case tm := <-c.TaskMsg:
				logrus.Infof("taskID:%s message:%s", c.TaskID, tm)
				ctx.AbortWithError(http.StatusBadRequest, errors.New(tm))
				ssEvent.Data = "CLOSE-WITH-ERROR:" + "taskID:" + c.TaskID + " " + "message" + tm
				ssEvent.Render(w)
				w.Flush()
			}
		}
	}()
}
