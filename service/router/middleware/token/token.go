package token

import (
	log "github.com/Sirupsen/logrus"
	"github.com/gin-gonic/gin"
)

func Refresh(c *gin.Context) {
	log.Info("refresh access token")
}
