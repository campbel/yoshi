package yoshi

import (
	"fmt"
	"reflect"
	"strings"
)

// options loads the values from arguments in the value v.
// it looks for the yoshi tag in the struct fields and loads
// the corresponding value from args
func options(v reflect.Value, args ...string) error {
	pargs := parseArgs(args, getBooleanFlags(v)...)
	fields := reflect.VisibleFields(v.Elem().Type())
LOOP:
	for _, parg := range pargs {
		// we've reached a command, stop loading options
		if parg.command != "" {
			return nil
		}
		for _, field := range fields {
			tags := strings.Split(field.Tag.Get(TagFlag), ",")
			for _, tag := range tags {
				if tag == parg.flag {
					if setter, ok := setterMap[field.Type.Kind()]; ok {
						val := v.Elem().FieldByName(field.Name)
						if err := setter(val, parg.value); err != nil {
							return err
						}
						continue LOOP
					} else {
						return fmt.Errorf("unsupported type: %s", field.Type.Kind())
					}
				}
			}
		}
		return fmt.Errorf("unknown flag: %s", parg.flag)
	}

	return nil
}

func defaults(v reflect.Value) error {
	fields := reflect.VisibleFields(v.Elem().Type())
	for _, field := range fields {
		tag := field.Tag.Get(TagDefault)
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

func getBooleanFlags(v reflect.Value) []string {
	var flags []string
	fields := reflect.VisibleFields(v.Elem().Type())
	for _, field := range fields {
		if field.Type.Kind() == reflect.Bool {
			flags = append(flags, strings.Split(field.Tag.Get(TagFlag), ",")...)
		}
	}
	return flags
}
