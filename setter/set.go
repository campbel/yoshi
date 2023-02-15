package setter

import (
	"errors"
	"fmt"
	"reflect"
	"strconv"
	"strings"
	"time"
)

func Set(val reflect.Value, value string) error {
	value, err := convertValue(val.Type(), value)
	if err != nil {
		return err
	}
	fn, ok := setterMap[val.Kind()]
	if !ok {
		return fmt.Errorf("unsupported type: %s", val.Kind())
	}
	return fn(val, value)
}

func Supports(kind reflect.Kind) bool {
	_, ok := setterMap[kind]
	return ok
}

func SetAny(target any, value string) error {
	val := reflect.ValueOf(target)
	if val.Kind() != reflect.Ptr {
		return errors.New("target must be a pointer")
	}
	return Set(val.Elem(), value)
}

// convertValue converts a string representation of a value to the correct type.
func convertValue(typ reflect.Type, value string) (string, error) {
	switch typ.String() {
	case "time.Duration":
		d, err := time.ParseDuration(value)
		if err != nil {
			return "", err
		}
		return strconv.Itoa(int(d)), nil
	}
	return value, nil
}

// setterMap standardizes setting values from strings.
// its important to have the same behavior for validation and execution.
var setterMap = map[reflect.Kind]func(reflect.Value, string) error{
	reflect.String: func(val reflect.Value, s string) error {
		val.SetString(s)
		return nil
	},
	reflect.Int: func(val reflect.Value, s string) error {
		v, err := strconv.Atoi(s)
		if err != nil {
			return err
		}
		val.SetInt(int64(v))
		return nil
	},
	reflect.Int64: func(val reflect.Value, s string) error {
		v, err := strconv.ParseInt(s, 10, 64)
		if err != nil {
			return err
		}
		val.SetInt(v)
		return nil
	},
	reflect.Bool: func(val reflect.Value, s string) error {
		if s == "" {
			s = "false"
		}
		v, err := strconv.ParseBool(s)
		if err != nil {
			return err
		}
		val.SetBool(v)
		return nil
	},
	reflect.Slice: func(val reflect.Value, s string) error {
		if val.IsNil() {
			val.Set(reflect.MakeSlice(val.Type(), 0, 0))
		}
		if s == "" {
			return nil
		}
		parts := strings.Split(s, ",")
		for _, part := range parts {
			switch val.Type().Elem().Kind() {
			case reflect.String:
				val.Set(reflect.Append(val, reflect.ValueOf(part)))
			case reflect.Int, reflect.Int64:
				v, err := strconv.Atoi(part)
				if err != nil {
					return err
				}
				val.Set(reflect.Append(val, reflect.ValueOf(v)))
			case reflect.Bool:
				v, err := strconv.ParseBool(part)
				if err != nil {
					return err
				}
				val.Set(reflect.Append(val, reflect.ValueOf(v)))
			}
		}
		return nil
	},
	reflect.Map: func(val reflect.Value, s string) error {
		if val.Type().Key().Kind() != reflect.String {
			return errors.New("map key must be string")
		}
		if val.IsNil() {
			val.Set(reflect.MakeMap(val.Type()))
		}
		if s == "" {
			return nil
		}
		parts := strings.Split(s, ",")
		for _, part := range parts {
			p := strings.Split(part, "=")
			if len(p) != 2 {
				return errors.New("map value must be in the form key=value")
			}
			switch val.Type().Elem().Kind() {
			case reflect.String:
				val.SetMapIndex(reflect.ValueOf(p[0]), reflect.ValueOf(p[1]))
			case reflect.Int, reflect.Int64:
				v, err := strconv.Atoi(p[1])
				if err != nil {
					return err
				}
				val.SetMapIndex(reflect.ValueOf(p[0]), reflect.ValueOf(v))
			case reflect.Bool:
				v, err := strconv.ParseBool(p[1])
				if err != nil {
					return err
				}
				val.SetMapIndex(reflect.ValueOf(p[0]), reflect.ValueOf(v))
			}
		}
		return nil
	},
}
