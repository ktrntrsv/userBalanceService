package api

import (
	"github.com/gin-gonic/gin"
	"github.com/ktrntrsv/userBalanceService/pkg/logger"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"net/http"
)

func NewRouter(handler *gin.Engine, l logger.Interface) {
	// Options
	handler.Use(gin.Logger())
	handler.Use(gin.Recovery())

	// Swagger
	swaggerHandler := ginSwagger.DisablingWrapHandler(swaggerFiles.Handler, "DISABLE_SWAGGER_HTTP_HANDLER")
	handler.GET("/swagger/*any", swaggerHandler)

	// K8s probe
	handler.GET("/healthz", func(c *gin.Context) {
		l.Info("Health")
		c.Status(http.StatusOK)
	})

	handler.PATCH("/account/:id/balance", enroll)
	handler.GET("/account/:id/balance", getBalance)
	handler.POST("/transaction/", startTransaction)
	handler.PATCH("/transaction/:id/abort", abortTransaction)
	handler.PATCH("/transaction/:id/approve", approveTransaction)
}
