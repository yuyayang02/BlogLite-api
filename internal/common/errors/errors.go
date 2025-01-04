package errors

import (
	"fmt"
)

type Error interface {
	error
	Code() int
	Message() string
}

type baseError struct {
	code    ErrorCode // 错误码
	message string    // 错误信息
}

func New(ec ErrorCode, msg string) Error {
	return &baseError{
		code:    ec,
		message: msg,
	}
}

func NewDomainError(msg string) Error {
	return &baseError{
		code:    ErrorCodeDomainError,
		message: msg,
	}
}

func NewDomainErrorf(format string, args ...any) Error {
	return &baseError{
		code:    ErrorCodeDomainError,
		message: fmt.Sprintf(format, args...),
	}
}

func Errorf(code ErrorCode, format string, args ...any) error {
	return &baseError{
		message: fmt.Sprintf(format, args...),
		code:    code,
	}
}

func (e *baseError) Error() string {
	return fmt.Sprintf("[%v]%s", e.code.Code, e.message)
}

func (e *baseError) Code() int {
	return e.code.Code
}

func (e *baseError) Message() string {
	return e.message
}

type WrapError interface {
	Error
	Cause() error // 原始错误
}

type wrapError struct {
	code    ErrorCode
	message string
	cause   error
}

func Wrap(err error, ec ErrorCode, msg string) error {
	if err == nil {
		return nil
	}

	return &wrapError{
		code:    ec,
		message: msg,
		cause:   err,
	}
}

func (e *wrapError) Error() string {
	return fmt.Sprintf("[%v]%s: %s", e.code.Code, e.message, e.cause.Error())
}

func (e *wrapError) Code() int {
	return e.code.Code
}

func (e *wrapError) Message() string {
	return e.message
}

func (e *wrapError) Cause() error {
	return e.cause
}

type causer interface {
	Cause() error
}

// Cause 找到错误链上最早出现的错误
func Cause(err error) error {
	for err != nil {
		cause, ok := err.(causer)
		if !ok {
			break
		}
		err = cause.Cause()
	}
	return err
}

// Reason 找到错误链上最早出现的满足common/errors.E rror接口的错误
func Reason(err error) Error {

	var lastError Error
	for err != nil {

		// 使用 e, ok := err.(Error), 而不是 errors.As(...)
		// 原因是：当前逻辑只需要判断接口方法而不需要对错误链进行检查，
		// 这种情况下使用err.(Error)满足当前逻辑且比errors.As(...)快约400倍
		if e, ok := err.(Error); ok {
			lastError = e
		}

		cause, ok := err.(causer)
		if !ok {
			break
		}
		err = cause.Cause()
	}

	return lastError
}
func InternalServiceError(err error) error {
	return Wrap(err, ErrorCodeInternalServiceError, "内部服务错误")
}

func InvalidParamsError(err error) error {
	return Wrap(err, ErrorCodeInvalidParamsError, "请求参数错误")
}
