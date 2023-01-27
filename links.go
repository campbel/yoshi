package yoshi

import (
	"fmt"
	"reflect"
)

// Link represents the order of commands for a given node tree
// and argument list. example:
// <app> foo bar baz
// link(node=foo) ➜ link(node=bar) ➜ link(node=baz)
type link struct {
	node *cmdNode
	next *link
	err  error
}

func (l *link) execute() {
	// Check method first, field second
	runMethod := l.node.cmdReference.MethodByName("Run")
	runField := l.node.cmdReference.Elem().FieldByName("Run")
	if runMethod.IsValid() && !runMethod.IsZero() {
		runMethod.Call([]reflect.Value{l.node.optionsReference.Elem()})
	} else if runField.IsValid() && !runField.IsZero() && !runField.IsNil() {
		l.node.runReference.Elem().Call([]reflect.Value{l.node.optionsReference.Elem()})
	}
	if l.next != nil {
		l.next.execute()
	}
}

func (l *link) help(usage ...string) string {
	if l.next == nil {
		return help(l.node.cmdReference.Elem().Type(), nil, append(usage, l.node.name)...)
	}
	return l.next.help(append(usage, l.node.name)...)
}

func buildLinks(name string, node *cmdNode, args args) *link {
	link := new(link)
	link.node = node
	for i, arg := range args {
		if arg.command != "" {
			command := node.commands.Get(arg.command)
			if command == nil {
				link.err = fmt.Errorf("unknown command %s", arg.command)
				return link
			}
			link.next = buildLinks(arg.command, command, args[i+1:])
			return link
		}
		if arg.flag != "" {
			for _, option := range node.options {
				for _, flag := range option.flags {
					if flag == arg.flag {
						if arg.value == "" {
							arg.value = option.def
						}
						setter, ok := setterMap[option.typ.Kind()]
						if !ok {
							link.err = fmt.Errorf("invalid type %s", option.typ.Kind().String())
							return link
						}
						err := setter(option.val, arg.value)
						if err != nil {
							link.err = err
							return link
						}
					}
				}
			}
		}
	}
	return link
}
