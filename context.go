package yoshi

import (
	"fmt"
	"os"
	"reflect"
)

type Context[T any] struct {
	App  *T
	root *Node
}

func Create[T any]() *Context[T] {
	var t T
	return &Context[T]{
		App:  &t,
		root: buildNodes(reflect.ValueOf(&t)),
	}
}

func (a *Context[T]) Run(args ...string) {
	a.run(os.Args[1:]...)
}

func (a *Context[T]) run(args ...string) {
	buildLinks(a.root, parseArgs(args)).execute()
}

func (a *Context[T]) Validate() error {
	var multiErr error
	errs := a.root.validate("root")
	for _, err := range errs {
		multiErr = fmt.Errorf("%s: %s", multiErr, err)
	}
	return multiErr
}
