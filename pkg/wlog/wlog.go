package wlog

import (
	"fmt"
	"runtime"
)

type Error struct {
	msg string
}

func errorWithStack(originErr error) Error {
	pc, codePath, codeLine, ok := runtime.Caller(2)
	if !ok {
		return Error{
			msg: fmt.Errorf("%v,and stack error", originErr).Error(),
		}
	}

	wrapErr := fmt.Sprintf("%s:%d:%v %v",
		codePath,
		codeLine,
		runtime.FuncForPC(pc).Name(),
		originErr,
	)
	return Error{
		msg: wrapErr,
	}
}

func Errorf(format string, a ...interface{}) Error {
	originErr := fmt.Errorf("'%v'", fmt.Sprintf(format, a...))
	return errorWithStack(originErr)
}

func WrapError(e error) Error {
	originErr := fmt.Errorf("\n  -%v", e.Error())
	return errorWithStack(originErr)
}

func (e Error) Error() string {
	return e.msg
}
