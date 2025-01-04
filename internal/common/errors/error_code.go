package errors

type ErrorCode struct {
	Code int
}

var (
	ErrorCodeDomainError = ErrorCode{1000}

	ErrorCodeResourceNotFound = ErrorCode{2001}
	ErrorCodeResourceIsExists = ErrorCode{2002}

	ErrorCodeInvalidParamsError = ErrorCode{4001}

	ErrorCodeInternalServiceError = ErrorCode{5000}
)
