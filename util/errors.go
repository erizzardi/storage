package util

type ResponseError struct {
	StatusCode int
	Err        error
}

func (e *ResponseError) Error() string {
	return e.Err.Error()
}
