package yoshi

import (
	"errors"
	"fmt"
	"reflect"
	"strconv"
	"strings"
)

func eval(t any) error {
	return evalT(reflect.ValueOf(t))
}

func evalT(val reflect.Value) error {
	fields := reflect.VisibleFields(val.Elem().Type())
	for _, structField := range fields {
		field := val.Elem().FieldByName(structField.Name)
		if field.Kind() == reflect.Struct {
			evalT(field.Addr())
			continue
		}
		tag := structField.Tag.Get(TagDefault)
		if tag == "" {
			continue
		}
		parser, ok := setterMap[field.Kind()]
		if !ok {
			return fmt.Errorf("invalid type %s", field.Kind().String())
		}
		err := parser(field, tag)
		if err != nil {
			return err
		}
	}
	return nil
}

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
		v, err := strconv.Atoi(s)
		if err != nil {
			return err
		}
		val.SetInt(int64(v))
		return nil
	},
	reflect.Bool: func(val reflect.Value, s string) error {
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
		parts := strings.Split(s, ",")
		for _, part := range parts {
			switch val.Type().Elem().Kind() {
			case reflect.String:
				val.Set(reflect.Append(val, reflect.ValueOf(part)))
			case reflect.Int, reflect.Int64:
				v, _ := strconv.Atoi(part)
				val.Set(reflect.Append(val, reflect.ValueOf(v)))
			case reflect.Bool:
				v, _ := strconv.ParseBool(part)
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
		parts := strings.Split(s, ",")
		for _, part := range parts {
			p := strings.Split(part, "=")
			if len(p) != 2 {
				continue
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
