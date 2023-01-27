package yoshi

import (
	"fmt"
	"reflect"
)

type link struct {
	Name    string
	Self    reflect.Value
	Options reflect.Value
	Run     reflect.Value
	Next    *link
}

func (l *link) execute() {
	run := l.Self.Elem().FieldByName("Run")
	if run.IsValid() && !run.IsZero() && !run.IsNil() {
		l.Run.Elem().Call([]reflect.Value{})
	}
	if l.Next != nil {
		l.Next.execute()
	}
}

func buildLinks(node *Node, args args) *link {
	link := new(link)
	link.Self = node.Value
	link.Run = node.Run
	link.Options = node.Opts
	for i, arg := range args {
		if arg.command != "" {
			link.Next = buildLinks(node.Commands[arg.command], args[i+1:])
			link.Next.Name = arg.command
			return link
		}
		if arg.flag != "" {
			for _, option := range node.Options {
				for _, flag := range option.Flags {
					if flag == arg.flag {
						if arg.value == "" {
							arg.value = option.Default
						}
						setter, ok := setterMap[option.Type.Kind()]
						if !ok {
							panic(fmt.Errorf("invalid type %s", option.Type.Kind().String()))
						}
						err := setter(option.Value, arg.value)
						if err != nil {
							panic(err)
						}
					}
				}
			}
		}
	}
	return link
}
