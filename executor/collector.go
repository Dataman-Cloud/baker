package executor

import (
	"github.com/Sirupsen/logrus"
	"github.com/gin-gonic/gin"
	"github.com/manucorporat/sse"
)

type Collector struct {
	WorkID    string
	TaskStats chan *TaskStats
}

type TaskStats struct {
	Code    int
	Message string
}

func NewCollector(workID string, taskStats chan *TaskStats) *Collector {
	return &Collector{
		WorkID:    workID,
		TaskStats: taskStats,
	}
}

func (c *Collector) Stream(ctx *gin.Context) chan bool {
	dst := make(chan bool)
	w := ctx.Writer
	close := w.CloseNotify()
	go func() {
		logrus.Info("Collector start")
		for {
			select {
			case <-close:
				logrus.Infof("Close Nodify") // nothing to do.
				dst <- true
				return
			case ts := <-c.TaskStats:
				var data string
				status := TaskStatusEnum[ts.Code]
				e := ts.Message
				if e != "" {
					logrus.Infof("workID:%s taskstatus:%s message:%s", c.WorkID, status, e)
					data = "workID:" + c.WorkID + " " + "status:" + status + " message:" + e
				} else {
					logrus.Infof("workID:%s taskstatus:%s", c.WorkID, status)
					data = "workID:" + c.WorkID + " " + "status:" + status
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
	return dst
}
