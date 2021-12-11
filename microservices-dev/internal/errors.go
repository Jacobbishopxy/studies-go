package internal

import "fmt"

// Error 结构体代表了一个错误可以包含其他错误，
// 其中还包含了一个错误码，用于区分被触发错误的不同方式。
type Error struct {
	orig error
	msg  string
	code ErrorCode
}

// 错误码
type ErrorCode uint

// 触发错误的枚举
const (
	ErrorCodeUnknown ErrorCode = iota
	ErrorCodeNotFound
	ErrorCodeInvalidArgument
)

// 返回一个被包裹的错误
func WrapErrorf(orig error, code ErrorCode, format string, a ...interface{}) error {
	return &Error{
		code: code,
		orig: orig,
		msg:  fmt.Sprintf(format, a...),
	}
}

// 新错误实例
func NewErrorf(code ErrorCode, format string, a ...interface{}) error {
	return WrapErrorf(nil, code, format, a...)
}

// Error 返回错误的字符串表示
func (e *Error) Error() string {
	if e.orig != nil {
		return fmt.Sprintf("%s: %s", e.msg, e.orig.Error())
	}
	return e.msg
}

// Error 返回原有错误
func (e *Error) Unwrap() error {
	return e.orig
}

// Error 返回错误码
func (e *Error) Code() ErrorCode {
	return e.code
}
