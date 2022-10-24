package api

import (
	"context"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/ktrntrsv/userBalanceService/pkg/logger"
	"net/http"
)

type accountRoutes struct {
	accountUsecase accountUsecase
	log            logger.Interface
}
type accountUsecase interface {
	EnrollBalance(ctx context.Context, accountId uuid.UUID, sum float64) error
	GetBalance(ctx context.Context, accountId uuid.UUID) (float64, error)
}

func newAccountHandlers(handler *gin.Engine, aUc accountUsecase, log logger.Interface) {
	r := &accountRoutes{accountUsecase: aUc, log: log}

	h := handler.Group("/account")
	{
		h.PATCH("/:id/balance", r.enroll)
		h.GET("/:id/balance", r.getBalance)
	}
}

func (r *accountRoutes) enroll(c *gin.Context) {

	type enrollRequest struct {
		Amount float64 `json:"amount" binding:"required"`
	}

	type enrollResponse struct {
		Status string `json:"status"`
	}

	accId, err := getAccountId(c)
	if err != nil {
		c.String(http.StatusNotFound, "wrong account id")
		return
	}

	var reqBody enrollRequest
	if err := c.Bind(&reqBody); err != nil {
		c.String(http.StatusBadRequest, "wrong parameters")
		return
	}

	if reqBody.Amount <= 0 {
		c.String(http.StatusBadRequest, "amount parameter is required")
		return
	}

	if err := r.accountUsecase.EnrollBalance(c, accId, reqBody.Amount); err != nil {
		c.Status(http.StatusInternalServerError)
		return
	}

	resp := enrollResponse{Status: "success"}
	c.JSON(http.StatusOK, resp)
}

func (r *accountRoutes) getBalance(c *gin.Context) {

	type getBalanceResponse struct {
		Balance float64 `json:"balance"`
	}

	accId, err := getAccountId(c)
	if err != nil {
		c.String(http.StatusNotFound, "wrong account id")
		return
	}
	balance, err := r.accountUsecase.GetBalance(c, accId)
	if err != nil {
		c.Status(http.StatusInternalServerError)
		return
	}
	resp := getBalanceResponse{balance}
	c.JSON(http.StatusOK, resp)
}

func getAccountId(c *gin.Context) (uuid.UUID, error) {
	id := c.Param("id")
	return uuid.Parse(id)
}
