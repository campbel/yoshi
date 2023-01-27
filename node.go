package yoshi

import (
	"fmt"
	"reflect"
	"strings"
)

type cmdNodes []*cmdNode

func (nodes cmdNodes) Get(name string) *cmdNode {
	for _, node := range nodes {
		if node.name == name {
			return node
		}
	}
	return nil
}

type cmdNode struct {
	name     string
	options  []cmdOption
	commands cmdNodes

	cmdReference     reflect.Value
	runReference     reflect.Value
	optionsReference reflect.Value
}

func (n cmdNode) validate(chain string) []error {
	var errs []error
	for _, option := range n.options {
		if err := option.validate(chain); err != nil {
			errs = append(errs, err...)
		}
	}
	for _, node := range n.commands {
		if err := node.validate(chain + "/" + node.name); err != nil {
			errs = append(errs, err...)
		}
	}
	return errs
}

func buildNodes(name string, val reflect.Value) *cmdNode {
	var node cmdNode = cmdNode{
		name:         name,
		commands:     make(cmdNodes, 0),
		cmdReference: val,
	}
	fields := reflect.VisibleFields(val.Elem().Type())
	for _, field := range fields {
		structField := val.Elem().FieldByName(field.Name)
		switch field.Name {
		case "Options":
			if structField.Kind() == reflect.Struct {
				node.options = buildNodeOptions(structField.Addr())
				node.optionsReference = structField.Addr()
			}
		case "Run":
			if structField.Kind() == reflect.Func {
				node.runReference = structField.Addr()
			}
		default:
			if structField.Kind() == reflect.Struct {
				node.commands = append(node.commands, buildNodes(strings.ToLower(field.Name), structField.Addr()))
			}
		}
	}
	return &node
}

// cmdOption represents a single option on a command.
type cmdOption struct {
	field string
	flags []string
	def   string
	desc  string

	typ reflect.Type
	val reflect.Value
}

func (o cmdOption) validate(chain string) []error {
	var errs []error
	if len(o.flags) == 0 {
		errs = append(errs, fmt.Errorf("missing flags for %s.%s", chain, o.field))
	}
	if o.desc == "" {
		errs = append(errs, fmt.Errorf("missing description for %s.%s", chain, o.field))
	}
	if o.def != "" {
		fn, ok := setterMap[o.typ.Kind()]
		if !ok {
			errs = append(errs, fmt.Errorf(`invalid type "%s" for %s.%s`, o.typ.Kind().String(), chain, o.field))
		}
		val := reflect.New(o.typ)
		err := fn(val.Elem(), o.def)
		if err != nil {
			errs = append(errs, fmt.Errorf(`invalid default value "%s" for %s.%s`, o.def, chain, o.field))
		}
	}
	return errs
}

func buildNodeOptions(val reflect.Value) []cmdOption {
	var options []cmdOption
	fields := reflect.VisibleFields(val.Elem().Type())
	for _, field := range fields {
		var option cmdOption
		option.field = field.Name
		option.flags = strings.Split(field.Tag.Get(TagFlag), ",")
		option.def = field.Tag.Get(TagDefault)
		option.typ = field.Type
		option.val = val.Elem().FieldByName(field.Name)
		if setter, ok := setterMap[option.typ.Kind()]; ok {
			setter(option.val, option.def)
		}
		option.desc = field.Tag.Get(TagDescription)
		options = append(options, option)
	}
	return options
}
