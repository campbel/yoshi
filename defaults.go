package yoshi

import (
	"reflect"
	"strconv"
	"strings"
)

func eval(t any) {
	evalT(reflect.ValueOf(t))
}

func evalT(val reflect.Value) {
	fields := reflect.VisibleFields(val.Elem().Type())
	for _, structField := range fields {
		field := val.Elem().FieldByName(structField.Name)
		if field.Kind() == reflect.Struct {
			evalT(field.Addr())
			continue
		}
		tag := structField.Tag.Get("yoshi-def")
		if tag == "" {
			continue
		}
		switch field.Kind() {
		case reflect.String:
			field.SetString(tag)
		case reflect.Int, reflect.Int64:
			v, _ := strconv.Atoi(tag)
			field.SetInt(int64(v))
		case reflect.Bool:
			v, _ := strconv.ParseBool(tag)
			field.SetBool(v)
		case reflect.Array, reflect.Slice:
			if field.IsNil() {
				field.Set(reflect.MakeSlice(field.Type(), 0, 0))
			}
			parts := strings.Split(tag, ",")
			for _, part := range parts {
				switch field.Type().Elem().Kind() {
				case reflect.String:
					field.Set(reflect.Append(field, reflect.ValueOf(part)))
				case reflect.Int, reflect.Int64:
					v, _ := strconv.Atoi(part)
					field.Set(reflect.Append(field, reflect.ValueOf(v)))
				case reflect.Bool:
					v, _ := strconv.ParseBool(part)
					field.Set(reflect.Append(field, reflect.ValueOf(v)))
				}
			}
		case reflect.Map:
			if field.Type().Key().Kind() != reflect.String {
				continue
			}
			if field.IsNil() {
				field.Set(reflect.MakeMap(field.Type()))
			}
			parts := strings.Split(tag, ",")
			for _, part := range parts {
				p := strings.Split(part, "=")
				if len(p) != 2 {
					continue
				}
				switch field.Type().Elem().Kind() {
				case reflect.String:
					field.SetMapIndex(reflect.ValueOf(p[0]), reflect.ValueOf(p[1]))
				case reflect.Int, reflect.Int64:
					v, _ := strconv.Atoi(p[1])
					field.SetMapIndex(reflect.ValueOf(p[0]), reflect.ValueOf(v))
				case reflect.Bool:
					v, _ := strconv.ParseBool(p[1])
					field.SetMapIndex(reflect.ValueOf(p[0]), reflect.ValueOf(v))
				}
			}
		}
	}
}
