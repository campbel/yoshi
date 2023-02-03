package options

import (
	"fmt"
	"reflect"
)

func CreateFromArgs(ts []reflect.Type, args []string) ([]reflect.Value, error) {
	var options []reflect.Value
	for _, t := range ts {
		val, err := createOption(t, args)
		if err != nil {
			return nil, err
		}
		options = append(options, val)
	}
	return options, nil
}

func createOption(t reflect.Type, args []string) (reflect.Value, error) {
	v := reflect.New(t)
	if err := defaults(v); err != nil {
		return v, err
	}
	if err := options(v, args); err != nil {
		return v, err
	}
	return v.Elem(), nil
}

func defaults(v reflect.Value) error {
	fields := reflect.VisibleFields(v.Elem().Type())
	for _, field := range fields {
		if field.Type.Kind() == reflect.Struct {
			if err := defaults(v.Elem().FieldByName(field.Name).Addr()); err != nil {
				return err
			}
			continue
		}
		tag := parseTags(field).Default
		if tag == "" {
			continue
		}
		setter, ok := setterMap[field.Type.Kind()]
		if !ok {
			return fmt.Errorf("unsupported type: %s", field.Type.Kind())
		}
		val := v.Elem().FieldByName(field.Name)
		if err := setter(val, tag); err != nil {
			return err
		}
	}
	return nil
}

// options loads the values from arguments in the value v.
// it looks for the yoshi tag in the struct fields and loads
// the corresponding value from args
func options(v reflect.Value, arguments []string) error {
	type address struct {
		v     reflect.Value
		field string
	}

	set := func(a address, value string) error {
		val := a.v.Elem().FieldByName(a.field)
		if setter, ok := setterMap[val.Kind()]; ok {
			if err := setter(val, value); err != nil {
				return err
			}
		} else {
			return fmt.Errorf("unsupported type: %s", val.Kind())
		}
		return nil
	}

	pos := make([]address, 0)
	dict := make(map[string]address)
	var lookup func(v reflect.Value)
	lookup = func(v reflect.Value) {
		for _, field := range reflect.VisibleFields(v.Elem().Type()) {
			if field.Type.Kind() == reflect.Struct {
				lookup(v.Elem().FieldByName(field.Name).Addr())
				continue
			}
			option := parseTags(field)
			if option.Positional() {
				pos = append(pos, address{v, field.Name})
				continue
			}
			for _, flag := range option.Flags {
				dict[flag] = address{v, field.Name}
			}
		}
	}
	lookup(v)

	positionalIndex := 0
	for _, arg := range parseArgs(arguments) {
		if arg.key == "" {
			if positionalIndex >= len(pos) {
				return fmt.Errorf("invalid argument: %s", arg.value)
			}
			if err := set(pos[positionalIndex], arg.value); err != nil {
				return err
			}
			positionalIndex++
			continue
		}
		if address, ok := dict[arg.key]; ok {
			if err := set(address, arg.value); err != nil {
				return err
			}
			continue
		}
		return fmt.Errorf("invalid flag: %s", arg.key)
	}

	return nil
}
