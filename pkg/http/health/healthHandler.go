package health

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

func Alive(c *gin.Context) {
	c.String(http.StatusOK, time.Now().String())
}
