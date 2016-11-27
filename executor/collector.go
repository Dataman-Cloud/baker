package executor

import (
	_ "errors"
	_ "net/http"

	"github.com/Sirupsen/logrus"
	"github.com/gin-gonic/gin"
	"github.com/manucorporat/sse"
)

type Collector struct {
	TaskID     string
	TaskStatus chan *TaskStatus
}

type TaskStatus struct {
	StatusCode int
	Message    string
}

func NewCollector(taskID string, taskStatus chan *TaskStatus) *Collector {
	return &Collector{
		TaskID:     taskID,
		TaskStatus: taskStatus,
	}
}

func (c *Collector) Stream(ctx *gin.Context) {
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
				var data string
				status := TaskStatusEnum[ts.StatusCode]
				logrus.Infof("taskID:%s status:%s message:%s", c.TaskID, status, ts.Message)
				data = "taskID:" + c.TaskID + " " + "status:" + status + "message:" + ts.Message
				sse.Encode(w, sse.Event{Event: "task-status", Data: data})
				if ts.StatusCode == StatusFinished || ts.StatusCode == StatusFailed || ts.StatusCode == StatusExpired {
					data = "CLOSE"
					sse.Encode(w, sse.Event{Event: "task-status", Data: data})
				}
			}
		}
	}()
}
