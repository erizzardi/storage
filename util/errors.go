package util

type ResponseError struct {
	StatusCode int    `json='statusCode'`
	Message    string `json='message'`
}

func (e *ResponseError) Error() string {
	return e.Message
}

//===========
// 400-errors
//===========
// 400 Bad request
type BadRequestError struct{}

func (BadRequestError) Error() string { return "Bad request" }

// 401 Unauthorized
type UnauthorizedError struct{}

func (UnauthorizedError) Error() string { return "Unauthorized" }

// 403 Forbidden
type ForbiddenError struct{}

func (ForbiddenError) Error() string { return "Forbidden" }

// 404 Not found
type NotFoundError struct{}

func (NotFoundError) Error() string { return "Not found" }

// 405 Method not allowed
type MethodNotAllowedError struct{}

func (MethodNotAllowedError) Error() string { return "Method not allowed" }

// 409 Conflict
type ConflictError struct{}

func (ConflictError) Error() string { return "Conflict" }

// 413 Payload too large
type PayloadTooLargeError struct{}

func (PayloadTooLargeError) Error() string { return "Payload too large" }

// 415 Unsupported media type
type UnsupportedMediaTypeError struct{}

func (UnsupportedMediaTypeError) Error() string { return "Unsupported media type" }

//===========
// 500-errors
//===========
// 500 Internal server error
type InternalServerError struct{}

func (InternalServerError) Error() string { return "Internal server error" }

// 504 Gateway timeout
type GatewayTimeoutError struct{}

func (GatewayTimeoutError) Error() string { return "Gateway timeout" }
