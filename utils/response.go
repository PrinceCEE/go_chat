package utils

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

type ServerError struct {
	Message    string
	StatusCode int
	Err        error
}

func (e *ServerError) Error() string {
	return e.Message
}
