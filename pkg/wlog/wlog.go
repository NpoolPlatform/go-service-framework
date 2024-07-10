package wlog

import (
	"errors"
	"fmt"
	"runtime"
)

type Error struct {
	originErr error
	msg       string
}

func errorWithStack(_wrapErr, originErr error) error {
	callStackNum := 2
	pc, codePath, codeLine, ok := runtime.Caller(callStackNum)
	if !ok {
		return &Error{
			msg: fmt.Errorf("%v,and stack error", originErr).Error(),
		}
	}

	wrapErr := fmt.Sprintf("%s:%d:%v %v",
		codePath,
		codeLine,
		runtime.FuncForPC(pc).Name(),
		_wrapErr,
	)
	return &Error{
		originErr: originErr,
		msg:       wrapErr,
	}
}

func Errorf(format string, a ...interface{}) error {
	originErr := fmt.Errorf(format, a...)
	return errorWithStack(originErr, originErr)
}

func WrapError(e error) error {
	if e == nil || e == (*Error)(nil) {
		return nil
	}
	wrapErr := fmt.Errorf("\n  -%w", e)
	originErr := e
	if _e, ok := e.(*Error); ok {
		originErr = _e.originErr
	}
	return errorWithStack(wrapErr, originErr)
}

func (e *Error) Error() string {
	return e.msg
}

func (e *Error) Unwrap() error {
	return e.originErr
}

func Equal(e1, e2 error) bool {
	if e1 == nil || e2 == nil {
		return e1 == e2
	}

	if _e1, ok := e1.(*Error); ok {
		e1 = _e1.originErr
	}

	if _e2, ok := e2.(*Error); ok {
		e2 = _e2.originErr
	}

	return unwrap(e1).Error() == unwrap(e2).Error()
}

func unwrap(err error) error {
	for {
		_err := errors.Unwrap(err)
		if _err == nil {
			return err
		}
		err = _err
	}
}

func Unwrap(err error) error {
	return unwrap(err)
}
