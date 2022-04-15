package view

import (
	"fmt"
	"net/http"
)

type ErrorCode int64

const (
	ErrorStatusBadRequest          ErrorCode = http.StatusBadRequest          // 400
	ErrorStatusUnauthorized        ErrorCode = http.StatusUnauthorized        // 401
	ErrorStatusPaymentRequired     ErrorCode = http.StatusPaymentRequired     // 402
	ErrorStatusForbidden           ErrorCode = http.StatusForbidden           // 403
	ErrorStatusNotFound            ErrorCode = http.StatusNotFound            // 404
	ErrorStatusInternalServerError ErrorCode = http.StatusInternalServerError // 500
)

type AppError struct {
	Code ErrorCode
	Msg  string
}

func (e AppError) Error() string {
	return fmt.Sprintf("[%d]%s", e.Code, e.Msg)
}

// NewBadRequestErrorFromModel 400
func NewBadRequestErrorFromModel(msg string) AppError {
	return AppError{
		Code: ErrorStatusBadRequest,
		Msg:  msg,
	}
}

// NewUnauthorizedErrorFromModel 401
func NewUnauthorizedErrorFromModel(msg string) AppError {
	return AppError{
		Code: ErrorStatusUnauthorized,
		Msg:  msg,
	}
}

func NewNotFoundErrorFromModel(msg string) AppError {
	return AppError{
		Code: ErrorStatusNotFound,
		Msg:  msg,
	}
}

func NewInternalServerErrorFromModel(msg string) AppError {
	return AppError{
		Code: ErrorStatusInternalServerError,
		Msg:  msg,
	}
}

func NewDBErrorFromModel(err error) AppError {
	if err.Error() == "sql: no rows in result set" {
		return NewNotFoundErrorFromModel("該当データがありません")
	}
	return NewInternalServerErrorFromModel(err.Error())
}
