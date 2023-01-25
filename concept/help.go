package concept

import (
	"bytes"
	"fmt"
	"reflect"
	"strings"
	"text/tabwriter"
)

func help(command reflect.Type, err error, commands ...string) string {
	output := ""
	if err != nil {
		output += "Error: " + err.Error() + "\n"
	}
	// commands
	if len(commands) > 0 {
		output += "Usage:"
		for _, cmd := range commands {
			output += " " + strings.ToLower(cmd)
		}
		output += " [options]\n"
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
				defaultValue := field.Tag.Get("yoshi-default")
				if defaultValue != "" {
					defaultValue = fmt.Sprintf(" (default: %s)", defaultValue)
				}
				fmt.Fprintf(w, "\n\t%s\t%s\t%s", field.Tag.Get("yoshi-flag"), field.Type.String(), field.Tag.Get("yoshi-desc")+defaultValue)
			}
			w.Flush()
			output += buffer.String() + "\n"
		}
	}

	return output
}
