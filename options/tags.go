package options

import (
	"reflect"
	"strings"
)

const (
	tagFlag = "yoshi"
)

type Option struct {
	Flags       []string
	Description string
	Default     string
	Type        string
}

func (o Option) Positional() bool {
	for _, flag := range o.Flags {
		if flag[0] == '-' {
			return false
		}
	}
	return true
}

func GetOptions(typ reflect.Type) []Option {
	var options []Option
	for _, field := range reflect.VisibleFields(typ) {
		kind := field.Type.Kind()
		if setterMap[kind] != nil {
			option := parseTags(field)
			if len(option.Flags) > 0 {
				options = append(options, option)
			}
		}
	}
	return options
}

func parseTags(field reflect.StructField) Option {
	option := Option{}
	tag := field.Tag.Get(tagFlag)
	if tag == "" {
		return option
	}
	option.Type = field.Type.String()
	parts := strings.Split(tag, ";")
	if len(parts) > 2 {
		option.Default = parts[2]
	}
	if len(parts) > 1 {
		option.Description = parts[1]
	}
	if len(parts) > 0 {
		if parts[0] == "" {
			option.Flags = []string{strings.ToUpper(field.Name)}
		} else {
			option.Flags = strings.Split(parts[0], ",")
		}
	}
	return option
}
