package executor

import (
	"errors"
	_ "net/http"

	"github.com/Sirupsen/logrus"
	"github.com/gin-gonic/gin"
	"github.com/manucorporat/sse"
)

type Collector struct {
	TaskID    string
	TaskStats chan *TaskStats
}

type TaskStats struct {
	Code    int
	Message string
}

func NewCollector(taskID string, taskStats chan *TaskStats) *Collector {
	return &Collector{
		TaskID:    taskID,
		TaskStats: taskStats,
	}
}

func (c *Collector) Stream(ctx *gin.Context) chan error {
	e := make(chan error)
	w := ctx.Writer
	clientClose := w.CloseNotify()
	go func() {
		logrus.Info("Collector start")
		for {
			select {
			case <-clientClose:
				logrus.Infof("Close Nodify") // nothing to do.
				e <- errors.New("abc")
				return
			case ts := <-c.TaskStats:
				var data string
				status := TaskStatusEnum[ts.Code]
				e := ts.Message
				if e != "" {
					logrus.Infof("taskID:%s status:%s message:%s", c.TaskID, status, e)
					data = "taskID:" + c.TaskID + " " + "status:" + status + " message:" + e
				} else {
					logrus.Infof("taskID:%s status:%s", c.TaskID, status)
					data = "taskID:" + c.TaskID + " " + "status:" + status
				}
				sse.Encode(w, sse.Event{Event: "task-status", Data: data})
				w.Flush()
				if ts.Code == StatusFinished || ts.Code == StatusFailed || ts.Code == StatusExpired {
					data = "CLOSE"
					sse.Encode(w, sse.Event{Event: "task-status", Data: data})
					w.Flush()
				}
			}
		}
	}()
	return e
}
