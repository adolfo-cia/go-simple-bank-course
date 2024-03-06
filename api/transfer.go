package api

import (
	"database/sql"
	"errors"
	"fmt"
	"net/http"

	db "github.com/adolfo-cia/go-simple-bank-course/db/sqlc"
	"github.com/adolfo-cia/go-simple-bank-course/token"
	"github.com/gin-gonic/gin"
)

type transferRequest struct {
	FromAccountID int64  `json:"fromAccountId" binding:"required,min=1"`
	ToAccountID   int64  `json:"toAccountId" binding:"required,min=1"`
	Amount        int64  `json:"amount" binding:"required,gt=0"`
	Currency      string `json:"currency" binding:"required,currency"`
}

func (s *Server) createTransfer(ctx *gin.Context) {
	var req transferRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	authPayload := ctx.MustGet(authorizationPayloadCtxKey).(*token.Payload)

	fromAccount, isValid := s.validAccount(ctx, req.FromAccountID, req.Currency)
	if !isValid {
		return
	}

	if fromAccount.Owner != authPayload.Username {
		ctx.JSON(
			http.StatusUnauthorized,
			errorResponse(errors.New("the 'from account' does not belong to authenticated user")))
		return
	}

	if _, isValid := s.validAccount(ctx, req.ToAccountID, req.Currency); !isValid {
		return
	}

	arg := db.TransferTxParam{
		FromAccountId: req.FromAccountID,
		ToAccountId:   req.ToAccountID,
		Amount:        req.Amount}

	result, err := s.store.TransferTx(ctx, arg)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusCreated, result)
}

func (s *Server) validAccount(ctx *gin.Context, accountId int64, currency string) (account db.Account, isValid bool) {
	account, err := s.store.GetAccount(ctx, accountId)
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, errorResponse(err))
		} else {
			ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		}
		return
	}

	if account.Currency != currency {
		err = fmt.Errorf("account [%d] currency mismatch: %s vs %s", accountId, account.Currency, currency)
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	isValid = true
	return
}
