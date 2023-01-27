package yoshi

import (
	"fmt"
	"os"
	"reflect"
)

type Context[T any] struct {
	App  *T
	name string
	root *Node
}

func Create[T any](name string) *Context[T] {
	var t T
	return &Context[T]{
		App:  &t,
		name: name,
		root: buildNodes(reflect.ValueOf(&t)),
	}
}

func (a *Context[T]) Run(args ...string) {
	a.run(os.Args[1:]...)
}

func (a *Context[T]) run(args ...string) {
	buildLinks(a.name, a.root, parseArgs(args)).execute()
}

func (a *Context[T]) Validate() []error {
	return a.root.validate(a.name)
}

func (a *Context[T]) PrintHelp() {
	fmt.Println(a.help(os.Args[1:]...))
}

func (a *Context[T]) help(args ...string) string {
	return buildLinks(a.name, a.root, parseArgs(args)).help()
}
