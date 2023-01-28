package yoshi

import (
	"fmt"
	"reflect"
)

// userError is an error that is caused by the user.
type userError struct {
	cmd reflect.Type
	err error
}

func (r *userError) Error() string {
	return r.err.Error()
}

func runErr(cmd reflect.Type, err error) *userError {
	return &userError{err: err}
}

func runErrf(cmd reflect.Type, format string, args ...interface{}) *userError {
	return &userError{cmd: cmd, err: fmt.Errorf(format, args...)}
}
