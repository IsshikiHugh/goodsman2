// Handlers about app service will be put here.

package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// Be used to check whether the service is online.
func Ping(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"err":  "NULL",
		"data": "Pong!",
	})
}
