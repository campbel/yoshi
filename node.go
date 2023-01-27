package yoshi

import (
	"fmt"
	"reflect"
	"strings"
)

type Node struct {
	Options  []Option
	Commands map[string]Node
}

func (n Node) Validate() []error {
	return n.validate("/root")
}

func (n Node) validate(chain string) []error {
	var errs []error
	for _, option := range n.Options {
		if err := option.validate(chain); err != nil {
			errs = append(errs, err...)
		}
	}
	for command, node := range n.Commands {
		if err := node.validate(chain + "/" + command); err != nil {
			errs = append(errs, err...)
		}
	}
	return errs
}

type Option struct {
	Field       string
	Flags       []string
	Type        reflect.Type
	Default     string
	Description string
}

func (o Option) validate(chain string) []error {
	var errs []error
	if len(o.Flags) == 0 {
		errs = append(errs, fmt.Errorf("missing flags for %s.%s", chain, o.Field))
	}
	if o.Description == "" {
		errs = append(errs, fmt.Errorf("missing description for %s.%s", chain, o.Field))
	}
	if o.Default != "" {
		fn, ok := setterMap[o.Type.Kind()]
		if !ok {
			errs = append(errs, fmt.Errorf(`invalid type "%s" for %s.%s`, o.Type.Kind().String(), chain, o.Field))
		}
		val := reflect.New(o.Type)
		err := fn(val, o.Default)
		if err != nil {
			errs = append(errs, fmt.Errorf(`invalid default value "%s" for %s.%s`, o.Default, chain, o.Field))
		}
	}
	return errs
}

func Describe[T any]() Node {
	var t T
	return describe(reflect.TypeOf(t))
}

func describe(typ reflect.Type) Node {
	var node Node = Node{Commands: make(map[string]Node)}
	fields := reflect.VisibleFields(typ)
	for _, field := range fields {
		switch field.Name {
		case "Options":
			node.Options = describeOptions(field.Type)
		default:
			node.Commands[field.Name] = describe(field.Type)
		}
	}
	return node
}

func describeOptions(typ reflect.Type) []Option {
	var options []Option
	fields := reflect.VisibleFields(typ)
	for _, field := range fields {
		var option Option
		option.Field = field.Name
		option.Flags = strings.Split(field.Tag.Get(TagFlag), ",")
		option.Default = field.Tag.Get(TagDefault)
		option.Type = field.Type
		option.Description = field.Tag.Get(TagDescription)
		options = append(options, option)
	}
	return options
}
