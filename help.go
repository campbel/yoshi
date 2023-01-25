package yoshi

import (
	"bytes"
	"fmt"
	"reflect"
	"regexp"
	"strings"
	"text/tabwriter"
)

var positionalRegex = regexp.MustCompile(`\[([0-9]+)\]`)

func Help[T any](cmds ...string) string {
	return HelpE[T](nil, cmds...)
}

func HelpE[T any](err error, cmds ...string) string {
	output := ""
	// error first
	if err != nil {
		output += "Error: " + err.Error() + "\n"
	}
	// then commands
	if len(cmds) > 0 {
		output += "Commands:\n"
		for _, cmd := range cmds {
			output += fmt.Sprintf("  %s\n", cmd)
		}
	}
	// then arguments
	output += argumentsText[T]()
	// then options
	output += optionsText[T]()
	return output
}

func argumentsText[T any]() string {
	var t T
	fields := reflect.VisibleFields(reflect.TypeOf(t))
	if fields == nil {
		return ""
	}
	var buffer bytes.Buffer
	w := tabwriter.NewWriter(&buffer, 0, 0, 1, ' ', 0)
	for _, field := range fields {
		tag := field.Tag.Get("opts")
		if positionalRegex.MatchString(tag) {
			fmt.Fprintf(w, "\n  %s\t%s\t%s", tag, field.Name, field.Type.String())
		}
	}
	w.Flush()
	if buffer.Len() == 0 {
		return ""
	}
	return "Arguments:" + buffer.String() + "\n"
}

func optionsText[T any]() string {
	var t T
	fields := reflect.VisibleFields(reflect.TypeOf(t))
	if fields == nil {
		return ""
	}
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
			continue
		}
		def := field.Tag.Get("default")
		if def != "" {
			def = "(default " + def + ")"
		}

		fmt.Fprintf(w, "\n  %s\t%s\t%s %s", strings.Join(vals, ", "), field.Type.String(), description, def)
	}
	w.Flush()
	if buffer.Len() == 0 {
		return ""
	}
	return "Options:" + buffer.String() + "\n"
}
