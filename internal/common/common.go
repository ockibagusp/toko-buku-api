package common

import "net/http"

// BaseResponseModel helpers
type BaseResponseModel struct {
	Status  int    `json:"status"`
	Message string `json:"message"`
	Data    any    `json:"data,omitempty"`
}

// custom response
func NewResponse[T any](status int, message string, data T) BaseResponseModel {
	return BaseResponseModel{
		Status:  status,
		Message: message,
		Data:    data,
	}
}

// returns http 200 OK
func StatusOK[T any](data T) BaseResponseModel {
	return BaseResponseModel{
		Status:  http.StatusOK,
		Message: "OK",
		Data:    data,
	}
}

// returns http 400
func StatusFail(message string) BaseResponseModel {
	return BaseResponseModel{
		Status:  http.StatusBadRequest,
		Message: message,
	}
}

// retuns http 401
func StatusUnauthorized(message string) BaseResponseModel {
	return BaseResponseModel{
		Status:  http.StatusUnauthorized,
		Message: message,
	}
}

// returns http 500
func UnhandledError() BaseResponseModel {
	return BaseResponseModel{
		Status:  http.StatusInternalServerError,
		Message: "Unhandled error occurred. Please try again later",
	}
}

// returns http 500
func StatusInternalServerError(message string) BaseResponseModel {
	return BaseResponseModel{
		Status:  http.StatusInternalServerError,
		Message: message,
	}
}

// returns http 404
func StatusNotFound(message string) BaseResponseModel {
	return BaseResponseModel{
		Status:  http.StatusNotFound,
		Message: message,
	}
}
