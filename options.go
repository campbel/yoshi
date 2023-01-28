package yoshi

import (
	"reflect"
	"strings"
)

// options loads the values from arguments in the value v.
// it looks for the yoshi tag in the struct fields and loads
// the corresponding value from args
func options(v reflect.Value, args ...string) error {
	pargs := parseArgs(args)
	fields := reflect.VisibleFields(v.Elem().Type())
	for _, parg := range pargs {
		// we've reached a command, stop loading options
		if parg.command != "" {
			return nil
		}
		for _, field := range fields {
			tags := strings.Split(field.Tag.Get("yoshi"), ",")
			for _, tag := range tags {
				if tag == parg.flag {
					if setter, ok := setterMap[field.Type.Kind()]; ok {
						val := v.Elem().FieldByName(field.Name)
						if err := setter(val, parg.value); err != nil {
							return err
						}
					}
				}
			}
		}
	}

	return nil
}
