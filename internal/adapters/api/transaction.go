package api

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func startTransaction(c *gin.Context) {
	c.Status(http.StatusOK)
}

func approveTransaction(c *gin.Context) {
	c.Status(http.StatusOK)
}

func abortTransaction(c *gin.Context) {
	c.Status(http.StatusOK)
}
