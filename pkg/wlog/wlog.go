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

func Equal(err, target error) bool {
	if target == nil || err == nil {
		return err == target
	}
	return unwarp(err).Error() == unwarp(target).Error()
}

func unwarp(err error) error {
	for {
		_err := errors.Unwrap(err)
		if _err == nil {
			return err
		}
		err = _err
	}
}
