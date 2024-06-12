package utils

import "errors"

var (
	ErrInternalServer  = errors.New("internal server error")
	ErrNotFound        = errors.New("not found")
	ErrBadRequest      = errors.New("bad request")
	ErrUnauthorized    = errors.New("unauthorized")
	ErrDuplicateRecord = errors.New("duplicate record")
)

type ResponseMeta struct {
	Page        int    `json:"page,omitempty"`
	TotalPages  int    `json:"total_pages,omitempty"`
	AccessToken string `json:"access_token,omitempty"`
}

type Response[T any] struct {
	Success bool          `json:"success"`
	Data    T             `json:"data"`
	Message string        `json:"message"`
	Meta    *ResponseMeta `json:"meta,omitempty"`
}

type ResponseGeneric Response[any]

type ServerError struct {
	Message    string
	StatusCode int
	Err        error
}

func (e *ServerError) Error() string {
	return e.Message
}
