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
	next    *link
}

func (l *link) execute() {
	// Check method first, field second
	runMethod := l.Self.MethodByName("Run")
	runField := l.Self.Elem().FieldByName("Run")
	if runMethod.IsValid() && !runMethod.IsZero() {
		runMethod.Call([]reflect.Value{l.Options.Elem()})
	} else if runField.IsValid() && !runField.IsZero() && !runField.IsNil() {
		l.Run.Elem().Call([]reflect.Value{l.Options.Elem()})
	}
	if l.next != nil {
		l.next.execute()
	}
}

func (l *link) help(usage ...string) string {
	if l.next == nil {
		return help(l.Self.Elem().Type(), nil, usage...)
	}
	return l.next.help(append(usage, l.Name)...)
}

func buildLinks(name string, node *Node, args args) *link {
	link := new(link)
	link.Name = name
	link.Self = node.Value
	link.Run = node.Run
	link.Options = node.Opts
	for i, arg := range args {
		if arg.command != "" {
			link.next = buildLinks(arg.command, node.Commands[arg.command], args[i+1:])
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
