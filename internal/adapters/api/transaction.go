package api

import (
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/ktrntrsv/userBalanceService/internal/domain/models"
	"github.com/ktrntrsv/userBalanceService/pkg/logger"
	"net/http"
)

type transactionRoutes struct {
	transactionUsecase transactionUsecase
	log                logger.Interface
}

type transactionUsecase interface {
	StartTransaction(ctx context.Context, dto models.TransactionStartDTO) (uuid.UUID, error)
	ApproveTransaction(ctx context.Context, transactID uuid.UUID) error
	AbortTransaction(ctx context.Context, transactID uuid.UUID) error
}

func newTransactionHandlers(handler *gin.Engine, tUc transactionUsecase, log logger.Interface) {
	r := &transactionRoutes{transactionUsecase: tUc, log: log}

	h := handler.Group("/transaction")
	{
		h.POST("/", r.startTransaction)
		h.PATCH("/:id/abort", r.abortTransaction)
		h.PATCH("/:id/approve", r.approveTransaction)
	}
}

func (r *transactionRoutes) startTransaction(c *gin.Context) {
	type startTransactionResponse struct {
		TransactionId uuid.UUID `json:"transactionId"`
	}

	var reqBody models.TransactionStartDTO
	if err := c.Bind(&reqBody); err != nil {
		c.String(http.StatusBadRequest, "wrong parameters")
		return
	}

	if reqBody.AccountToId == reqBody.AccountFromId {
		c.Status(http.StatusBadRequest)
		return
	}

	id, err := r.transactionUsecase.StartTransaction(c, reqBody)
	if err != nil {
		fmt.Println("err", err)
		c.String(http.StatusInternalServerError, "can not create transaction")
		return
	}
	resp := startTransactionResponse{TransactionId: id}
	c.JSON(http.StatusOK, resp)
}

func (r *transactionRoutes) approveTransaction(c *gin.Context) {
	trId, err := getTransactionId(c)
	if err != nil {
		c.String(http.StatusNotFound, "can not find a transaction")
		return
	}

	err = r.transactionUsecase.ApproveTransaction(c, trId)
	if err != nil {
		c.String(http.StatusInternalServerError, "can not approve a transaction")
		return
	}
	c.Status(http.StatusOK)

}

func (r *transactionRoutes) abortTransaction(c *gin.Context) {
	trId, err := getTransactionId(c)
	if err != nil {
		c.String(http.StatusNotFound, "can not find a transaction")
		return
	}

	err = r.transactionUsecase.AbortTransaction(c, trId)
	if err != nil {
		c.String(http.StatusInternalServerError, "can not abort a transaction. (Аборт это убийство)")
		return
	}
	c.Status(http.StatusOK)

}

func getTransactionId(c *gin.Context) (uuid.UUID, error) {
	trId := c.Param("id")
	return uuid.Parse(trId)
}
