package yoshi

import (
	"bytes"
	"fmt"
	"reflect"
	"strings"
	"text/tabwriter"
)

func help(typ reflect.Type, err error, usage ...string) string {
	commands, options := getFields(typ)
	output := ""
	if err != nil && err != errHelp {
		output += "Error: " + err.Error() + "\n"
	}
	// usage
	if len(usage) > 0 {
		output += "Usage:"
		for _, cmd := range usage {
			output += " " + strings.ToLower(cmd)
		}
		if len(commands) > 0 {
			output += " COMMAND"
		} else {
			output += " [OPTIONS]"
		}
		output += "\n"
	}
	// commands
	if len(commands) > 0 {
		output += "Commands:\n"
		for _, cmd := range commands {
			output += "  " + strings.ToLower(cmd) + "\n"
		}
	}
	// options
	if len(options) > 0 {
		output += "Options:"
		var buffer bytes.Buffer
		w := tabwriter.NewWriter(&buffer, 0, 0, 1, ' ', 0)
		for _, opt := range options {
			field, ok := typ.FieldByName(opt)
			if !ok {
				panic("field not found: " + opt)
			}
			// tag
			tag := field.Tag.Get(TagFlag)
			if tag == "" {
				continue
			}
			line := "\n  " + tag
			// type
			typ := field.Type.String()
			if typ != "" {
				line += "\t" + typ
			}
			// description
			description := field.Tag.Get(TagDescription)
			if description != "" {
				line += "\t" + description
			}
			// default
			defaultValue := field.Tag.Get(TagDefault)
			if defaultValue != "" {
				line += fmt.Sprintf(" (default: %s)", defaultValue)
			}
			fmt.Fprint(w, line)
		}
		w.Flush()
		output += buffer.String() + "\n"
	}
	return output
}

func getFields(typ reflect.Type) ([]string, []string) {
	var commands []string
	var options []string
	for _, field := range reflect.VisibleFields(typ) {
		kind := field.Type.Kind()
		if kind == reflect.Func || (kind == reflect.Struct && !field.Anonymous) {
			commands = append(commands, field.Name)
		} else if setterMap[kind] != nil {
			options = append(options, field.Name)
		}
	}
	return commands, options
}
