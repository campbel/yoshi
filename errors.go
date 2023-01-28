package yoshi

import (
	"errors"
	"fmt"
	"reflect"
)

// userError is an error that is caused by the user.
type userError struct {
	loc reflect.Type
	err error
}

func (r *userError) Error() string {
	return r.err.Error()
}

func userErr(loc reflect.Type, err error) *userError {
	return &userError{loc: loc, err: err}
}

func userErrf(loc reflect.Type, format string, args ...interface{}) *userError {
	return &userError{loc: loc, err: fmt.Errorf(format, args...)}
}

var errHelp = errors.New("help")
