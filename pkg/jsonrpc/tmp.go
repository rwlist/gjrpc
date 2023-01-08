package jsonrpc

// TODO: remove the use of Result in router and delete this file.
type Result interface{}

// ErrorConverter is for converting error to jsonrpc error.
// TODO: think of error -> ErrorContext, with router path, etc.
type ErrorConverter func(error) (Result, *Error)

func DefaultErrorConverter(err error) (Result, *Error) {
	return nil, UnknownError.WithError(err)
}
