package opts

import (
	"bytes"
	"fmt"
	"reflect"
	"regexp"
	"strings"
	"text/tabwriter"
)

var positionalRegex = regexp.MustCompile(`\[([0-9]+)\]`)

func Help[T any](err ...error) string {
	var t T
	positionals := []string{}
	fields := reflect.VisibleFields(reflect.TypeOf(t))
	var buffer bytes.Buffer
	w := tabwriter.NewWriter(&buffer, 0, 0, 1, ' ', 0)
	for _, field := range fields {
		tag := field.Tag.Get("opts")
		if tag == "" {
			continue
		}
		description := field.Tag.Get("desc")
		vals := strings.Split(tag, ",")
		if positionalRegex.MatchString(vals[0]) {
			positionals = append(positionals, strings.ToUpper(field.Name))
			continue
		}
		def := field.Tag.Get("default")
		if def != "" {
			def = "(default " + def + ")"
		}

		fmt.Fprintf(w, "  %s\t%s\t%s %s\n", strings.Join(vals, ", "), field.Type.String(), description, def)
	}
	w.Flush()
	output := ""
	if len(err) > 0 {
		output += "Error: " + err[0].Error() + "\n\n"
	}
	output += "\n"
	output += "Options:\n"
	output += buffer.String()
	return output
}
