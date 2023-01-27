package yoshi

import (
	"fmt"
	"reflect"
	"strings"
)

type Node struct {
	Options  []NodeOption
	Commands map[string]*Node

	Type  reflect.Type
	Value reflect.Value
	Run   reflect.Value
	Opts  reflect.Value
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

func buildNodes(val reflect.Value) *Node {
	var node Node = Node{
		Commands: make(map[string]*Node),
		Value:    val,
		Type:     val.Elem().Type(),
	}
	fields := reflect.VisibleFields(node.Type)
	for _, field := range fields {
		structField := val.Elem().FieldByName(field.Name)
		switch field.Name {
		case "Options":
			node.Options = buildNodeOptions(structField.Addr())
			node.Opts = structField.Addr()
		case "Run":
			if structField.Kind() != reflect.Func {
				panic(fmt.Errorf("invalid type %s for Run", structField.Kind().String()))
			}
			node.Run = structField.Addr()
		default:
			if structField.Kind() != reflect.Struct {
				continue
			}
			node.Commands[strings.ToLower(field.Name)] = buildNodes(structField.Addr())
		}
	}
	return &node
}

type NodeOption struct {
	Field       string
	Flags       []string
	Default     string
	Description string

	Type  reflect.Type
	Value reflect.Value
}

func (o NodeOption) validate(chain string) []error {
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

func buildNodeOptions(val reflect.Value) []NodeOption {
	var options []NodeOption
	fields := reflect.VisibleFields(val.Elem().Type())
	for _, field := range fields {
		var option NodeOption
		option.Field = field.Name
		option.Flags = strings.Split(field.Tag.Get(TagFlag), ",")
		option.Default = field.Tag.Get(TagDefault)
		option.Type = field.Type
		option.Value = val.Elem().FieldByName(field.Name)
		option.Description = field.Tag.Get(TagDescription)
		options = append(options, option)
	}
	return options
}
