package utils

import "net/http"

// BaseResponseModel helpers
type BaseResponseModel[T any] struct {
	Status  int    `json:"status"`
	Message string `json:"message"`
}

// BaseResponseDataModel helpers
type BaseResponseDataModel[T any] struct {
	BaseResponseModel[T]
	Data T `json:"data,omitempty"`
}

// BaseResponseErrorModel helpers
type BaseResponseErrorModel[T any] struct {
	BaseResponseModel[T]
	Error T `json:"error,omitempty"`
}

// custom response
func NewResponse[T any](status int, message string, data T) BaseResponseDataModel[T] {
	return BaseResponseDataModel[T]{
		BaseResponseModel: BaseResponseModel[T]{
			Status:  status,
			Message: message,
		},
		Data: data,
	}
}

// custom response error
func NewResponseError[T any](status int, message string, err T) BaseResponseErrorModel[T] {
	return BaseResponseErrorModel[T]{
		BaseResponseModel: BaseResponseModel[T]{
			Status:  status,
			Message: message,
		},
		Error: err,
	}
}

// returns http 200 OK
func StatusOK[T any](data T) BaseResponseDataModel[T] {
	return BaseResponseDataModel[T]{
		BaseResponseModel: BaseResponseModel[T]{
			Status:  http.StatusOK,
			Message: "OK",
		},
		Data: data,
	}
}

// returns http 400
func StatusBadRequest[T string](message string) BaseResponseModel[T] {
	return BaseResponseModel[T]{
		Status:  http.StatusBadRequest,
		Message: message,
	}
}

// retuns http 401
func StatusUnauthorized[T string](message string) BaseResponseModel[T] {
	return BaseResponseModel[T]{
		Status:  http.StatusUnauthorized,
		Message: message,
	}
}

// returns http 500
func UnhandledError[T any]() BaseResponseModel[T] {
	return BaseResponseModel[T]{
		Status:  http.StatusInternalServerError,
		Message: "Unhandled error occurred. Please try again later",
	}
}

// returns http 500
func StatusInternalServerError[T string](message string) BaseResponseModel[T] {
	return BaseResponseModel[T]{
		Status:  http.StatusInternalServerError,
		Message: message,
	}
}

// returns http 404
func StatusNotFound[T string](message string) BaseResponseModel[T] {
	return BaseResponseModel[T]{
		Status:  http.StatusNotFound,
		Message: message,
	}
}
