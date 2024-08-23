package errors

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type AppError struct {
	Message    string `json:"message"`
	StatusCode int    `json:"status_code"`
	Err        error  `json:"-"`
}

func (e AppError) Error() string {
	if e.Err != nil {
		return fmt.Sprintf("%s: %v", e.Message, e.Err)
	}
	return e.Message
}

func NewAppError(message string, statusCode int, err error) *AppError {
	return &AppError{
		Message:    message,
		StatusCode: statusCode,
		Err:        err,
	}
}

func HandleError(w http.ResponseWriter, appErr *AppError) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(appErr.StatusCode)
	json.NewEncoder(w).Encode(appErr)
}
