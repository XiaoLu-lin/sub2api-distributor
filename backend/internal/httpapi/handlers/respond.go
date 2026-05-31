package handlers

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/lhl/sub2api-distributor/backend/internal/auth"
	"github.com/lhl/sub2api-distributor/backend/internal/distributor"
)

const internalServerMessage = "系统开小差了，请稍后重试"

func operatorMutationStatus(err error) int {
	switch {
	case errors.Is(err, distributor.ErrWithdrawalNotFound):
		return http.StatusNotFound
	case errors.Is(err, distributor.ErrInvalidWithdrawalTransition):
		return http.StatusBadRequest
	default:
		return http.StatusInternalServerError
	}
}

func respondError(c *gin.Context, status int, err error) {
	message := messageForError(status, err)
	c.JSON(status, gin.H{"message": message})
}

func messageForError(status int, err error) string {
	if err == nil {
		if status >= http.StatusInternalServerError {
			return internalServerMessage
		}
		return "请求失败"
	}

	switch {
	case errors.Is(err, auth.ErrInvalidCredentials),
		errors.Is(err, auth.ErrInactiveUser),
		errors.Is(err, auth.ErrDistributorNotEnabled),
		errors.Is(err, auth.ErrInvalidToken),
		errors.Is(err, distributor.ErrProfileNotFound),
		errors.Is(err, distributor.ErrInviteMetaNotFound),
		errors.Is(err, distributor.ErrInvalidProfile),
		errors.Is(err, distributor.ErrInvalidWithdrawalAmount),
		errors.Is(err, distributor.ErrWithdrawalAmountTooLarge),
		errors.Is(err, distributor.ErrWithdrawalNotFound),
		errors.Is(err, distributor.ErrInvalidWithdrawalTransition):
		return err.Error()
	default:
		if status >= http.StatusInternalServerError {
			return internalServerMessage
		}
		return err.Error()
	}
}
