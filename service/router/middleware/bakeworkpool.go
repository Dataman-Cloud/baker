package middleware

import (
	"strconv"

	"github.com/Sirupsen/logrus"
	"github.com/gin-gonic/gin"

	"github.com/Dataman-Cloud/baker/config"
	"github.com/Dataman-Cloud/baker/executor"
)

// BakeWorkPool is a middleware function that initializes the BakeWorkPool and attaches to
// the context of every http.Request.
func BakeWorkPool(cf *config.Config) gin.HandlerFunc {
	v := setupBakeWorkPool(cf)
	return func(c *gin.Context) {
		c.Set("bakeworkpool", v)
	}
}

// helper function to create the bakeWorkPool from the CLI context.
func setupBakeWorkPool(cf *config.Config) *executor.WorkPool {
	workPool := cf.WorkPool
	imagePushWorkPoolOptions := workPool["imagepush"]
	imagePushWorkPoolSize, _ := strconv.Atoi(imagePushWorkPoolOptions.MaxWorkers)
	if imagePushWorkPoolSize < 1 {
		logrus.Errorf("must provide positive maxWorkers; provided %d", imagePushWorkPoolSize)
		return nil
	}
	pool, err := executor.NewWorkPool(imagePushWorkPoolSize)
	if err != nil {
		logrus.Errorf("error create bakeworkpool: %s", err)
		return nil
	}
	return pool
}
