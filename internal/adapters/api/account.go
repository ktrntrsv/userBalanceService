package api

import (
	"context"
	"github.com/gin-gonic/gin"
	"github.com/ktrntrsv/userBalanceService/pkg/logger"
	"net/http"
)

type accountRoutes struct {
	accountUsecase accountUsecase
	log            logger.Interface
}
type accountUsecase interface {
	EnrollBalance(ctx context.Context, accountId string, sum float64) error
	GetBalance(ctx context.Context, accountId string) (float64, error)
}

func newAccountHandlers(handler *gin.Engine, aUc accountUsecase, l logger.Interface) {
	r := &accountRoutes{accountUsecase: aUc, log: l}

	h := handler.Group("/account")
	{
		h.PATCH("/:id/balance", r.enroll)
		h.GET("/:id/balance", r.getBalance)
	}
}

type enrollRequest struct {
	Amount float64 `json:"amount"`
}

type enrollResponse struct {
	Status string `json:"status"`
}

func (r *accountRoutes) enroll(c *gin.Context) {
	accId := getAccountId(c)

	var reqBody enrollRequest
	if err := c.Bind(reqBody); err != nil {
		c.Status(http.StatusBadRequest)
		return
	}

	if err := r.accountUsecase.EnrollBalance(c, accId, reqBody.Amount); err != nil {
		c.Status(http.StatusInternalServerError) // TODO: обработать ошибку
		return
	}

	resp := enrollResponse{Status: "success"}
	c.JSON(http.StatusOK, resp)
}

func (r *accountRoutes) getBalance(c *gin.Context) {
	//accId := getAccountId(c)

}

func getAccountId(c *gin.Context) string {
	return c.Param("id")
}
