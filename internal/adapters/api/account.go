package api

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func enroll(c *gin.Context) {
	c.Status(http.StatusOK)
}

func getBalance(c *gin.Context) {
	uid := c.Param("user_id")
	c.Status(http.StatusOK)
	_ = uid
}
