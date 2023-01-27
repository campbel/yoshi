package yoshi

import (
	"fmt"
	"reflect"
)

type link struct {
	name    string
	self    reflect.Value
	options reflect.Value
	run     reflect.Value
	error   error
	next    *link
}

func (l *link) execute() {
	// Check method first, field second
	runMethod := l.self.MethodByName("Run")
	runField := l.self.Elem().FieldByName("Run")
	if runMethod.IsValid() && !runMethod.IsZero() {
		runMethod.Call([]reflect.Value{l.options.Elem()})
	} else if runField.IsValid() && !runField.IsZero() && !runField.IsNil() {
		l.run.Elem().Call([]reflect.Value{l.options.Elem()})
	}
	if l.next != nil {
		l.next.execute()
	}
}

func (l *link) help(usage ...string) string {
	if l.next == nil {
		return help(l.self.Elem().Type(), nil, append(usage, l.name)...)
	}
	return l.next.help(append(usage, l.name)...)
}

func buildLinks(name string, node *Node, args args) *link {
	link := new(link)
	link.name = name
	link.self = node.Value
	link.run = node.Run
	link.options = node.Opts
	for i, arg := range args {
		if arg.command != "" {
			command, ok := node.Commands[arg.command]
			if !ok {
				link.error = fmt.Errorf("unknown command %s", arg.command)
				return link
			}
			link.next = buildLinks(arg.command, command, args[i+1:])
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
							link.error = fmt.Errorf("invalid type %s", option.Type.Kind().String())
							return link
						}
						err := setter(option.Value, arg.value)
						if err != nil {
							link.error = err
							return link
						}
					}
				}
			}
		}
	}
	return link
}
