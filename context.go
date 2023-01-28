package yoshi

import (
	"fmt"
	"os"
	"reflect"
)

type Context[T any] struct {
	App  *T
	name string
	root *cmdNode
}

func Create[T any](name string) *Context[T] {
	var t T
	return &Context[T]{
		App:  &t,
		name: name,
		root: buildNodes(name, reflect.ValueOf(&t)),
	}
}

func (a *Context[T]) Run(args ...string) {
	if len(args) == 0 {
		args = os.Args
	}
	a.run(args[1:]...)
}

func (a *Context[T]) run(args ...string) {
	buildLinks(a.name, a.root, parseArgs(args)).execute()
}

func (a *Context[T]) Validate() []error {
	return a.root.validate(a.name)
}

func (a *Context[T]) PrintHelp() {
	fmt.Println(a.Help(os.Args[1:]...))
}

func (a *Context[T]) Help(args ...string) string {
	return buildLinks(a.name, a.root, parseArgs(args)).help()
}
