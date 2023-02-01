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

// options loads the values from arguments in the value v.
// it looks for the yoshi tag in the struct fields and loads
// the corresponding value from args
func options(v reflect.Value, args []string) error {
	pargs := parseArgs(args)
	positionals, flags, err := parseOptions(v)
	if err != nil {
		return err
	}
	positionalIndex := 0
	for _, parg := range pargs {
		if parg.key == "" {
			if positionalIndex >= len(positionals) {
				return fmt.Errorf("invalid argument: %s", parg.value)
			}
			fieldName := positionals[positionalIndex]
			val := v.Elem().FieldByName(fieldName)
			if setter, ok := setterMap[val.Kind()]; ok {
				if err := setter(val, parg.value); err != nil {
					return err
				}
			} else {
				return fmt.Errorf("unsupported type: %s", val.Kind())
			}
			positionalIndex++
			continue
		}
		if parg.key != "" {
			if fieldName, ok := flags[parg.key]; ok {
				val := v.Elem().FieldByName(fieldName)
				if setter, ok := setterMap[val.Kind()]; ok {
					if err := setter(val, parg.value); err != nil {
						return err
					}
				} else {
					return fmt.Errorf("unsupported type: %s", val.Kind())
				}
				continue
			}
			return fmt.Errorf("invalid flag: %s", parg.key)
		}
	}
	return nil
}

func parseOptions(v reflect.Value) ([]string, map[string]string, error) {
	positionals := make([]string, 0)
	fields := reflect.VisibleFields(v.Elem().Type())
	flagMap := make(map[string]string)
	for _, field := range fields {
		flags := parseOption(field).Flags
		for _, flag := range flags {
			// Ignore empty flags
			if flag == "" {
				continue
			}
			if flag[0] == '-' {
				// handle flag
				if _, ok := flagMap[flag]; ok {
					return nil, nil, fmt.Errorf("duplicate flag: %s", flag)
				}
				flagMap[flag] = field.Name
			} else {
				// handle positional
				positionals = append(positionals, field.Name)
			}
		}
	}
	return positionals, flagMap, nil
}

func defaults(v reflect.Value) error {
	fields := reflect.VisibleFields(v.Elem().Type())
	for _, field := range fields {
		tag := parseOption(field).Default
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
