package yoshi

import (
	"bytes"
	"fmt"
	"reflect"
	"strings"
	"text/tabwriter"
)

func help(command reflect.Type, err error, usage ...string) string {
	subCommands := getSubCommands(command)
	output := ""
	if err != nil {
		output += "Error: " + err.Error() + "\n"
	}
	// usage
	if len(usage) > 0 {
		output += "Usage:"
		for _, cmd := range usage {
			output += " " + strings.ToLower(cmd)
		}
		output += " [options]"
		if len(subCommands) > 0 {
			if _, hasRun := command.FieldByName("Run"); hasRun {
				output += " [COMMAND]"
			} else {
				output += " COMMAND"
			}
		}
		output += "\n"
	}
	// commands
	if len(subCommands) > 0 {
		output += "Commands:\n"
		for _, cmd := range subCommands {
			output += "  " + strings.ToLower(cmd) + "\n"
		}
	}
	// options
	field, ok := command.FieldByName("Options")
	if ok {
		fields := reflect.VisibleFields(field.Type)
		if len(fields) > 0 {
			output += "Options:"
			var buffer bytes.Buffer
			w := tabwriter.NewWriter(&buffer, 0, 0, 1, ' ', 0)
			for _, field := range fields {
				line := ""
				line += "\n"
				// tag
				tag := field.Tag.Get(TagFlag)
				if tag != "" {
					line += "  " + tag
				}
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
	}

	return output
}

func getSubCommands(command reflect.Type) []string {
	var subCommands []string
	for _, field := range reflect.VisibleFields(command) {
		if field.Name == "Options" || field.Name == "Run" {
			continue
		}
		if field.Type.Kind() != reflect.Struct {
			continue
		}
		subCommands = append(subCommands, field.Name)
	}
	return subCommands
}
