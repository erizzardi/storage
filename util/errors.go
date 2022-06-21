package util

import "reflect"

type ResponseError struct {
	StatusCode int    `json:"statusCode"`
	Message    string `json:"message"`
}

func (e *ResponseError) Error() string {
	return e.Message
}

//============
// HTTP ERRORS
// ===========

//===========
// 4xx errors
//===========

// 400 Bad request
type BadRequestError struct{ Message string }

func (e BadRequestError) Error() string { return "Bad request: " + e.Message }

// type BadRequestError

// 401 Unauthorized
type UnauthorizedError struct{ Message string }

func (e UnauthorizedError) Error() string { return "Unauthorized: " + e.Message }

// 403 Forbidden
type ForbiddenError struct{ Message string }

func (e ForbiddenError) Error() string { return "Forbidden: " + e.Message }

// 404 Not found
type NotFoundError struct{ Message string }

func (e NotFoundError) Error() string { return "Not found: " + e.Message }

// 405 Method not allowed
type MethodNotAllowedError struct{ Message string }

func (e MethodNotAllowedError) Error() string { return "Method not allowed: " + e.Message }

// 409 Conflict
type ConflictError struct{ Message string }

func (e ConflictError) Error() string { return "Conflict : " + e.Message }

// 413 Payload too large
type PayloadTooLargeError struct{ Message string }

func (e PayloadTooLargeError) Error() string { return "Payload too large: " + e.Message }

// 415 Unsupported media type
type UnsupportedMediaTypeError struct{ Message string }

func (e UnsupportedMediaTypeError) Error() string { return "Unsupported media type: " + e.Message }

//===========
// 5xx errors
//===========

// 500 Internal server error
type InternalServerError struct{ Message string }

func (e InternalServerError) Error() string { return "Internal server error: " + e.Message }

// 504 Gateway timeout
type GatewayTimeoutError struct{ Message string }

func (e GatewayTimeoutError) Error() string { return "Gateway timeout: " + e.Message }

//============
// Miscellanea
//============
func ErrorIs(err error, target error) bool {
	return reflect.TypeOf(err) == reflect.TypeOf(target)
}
