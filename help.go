package yoshi

import (
	"bytes"
	"fmt"
	"reflect"
	"strings"
	"text/tabwriter"
)

func help(typ reflect.Type, err error, usage ...string) string {
	commands, positionals, options := getFields(typ)
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
		for _, pos := range positionals {
			output += " " + strings.ToUpper(pos)
		}
		if len(commands) > 0 {
			output += " COMMAND"
		}
		if len(options) > 0 {
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
			tags := getTags(field).Flags
			if len(tags) == 0 {
				continue
			}
			line := "\n  " + strings.Join(tags, ",")
			// type
			typ := field.Type.String()
			if typ != "" {
				line += "\t" + typ
			}
			// description
			description := getTags(field).Description
			if description != "" {
				line += "\t" + description
			}
			// default
			defaultValue := getTags(field).Default
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

func getFields(typ reflect.Type) ([]string, []string, []string) {
	var commands []string
	var positionals []string
	var options []string
	for _, field := range reflect.VisibleFields(typ) {
		kind := field.Type.Kind()
		if kind == reflect.Func || (kind == reflect.Struct && !field.Anonymous) {
			commands = append(commands, field.Name)
		} else if setterMap[kind] != nil {
			tags := getTags(field).Flags
			if len(tags) == 0 {
				continue
			}
			if tags[0][0] == '-' {
				options = append(options, field.Name)
			} else {
				positionals = append(positionals, field.Name)
			}
		}
	}
	return commands, positionals, options
}
