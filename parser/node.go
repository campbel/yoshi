package parser

import (
	"bytes"
	"fmt"
	"reflect"
	"strings"
	"text/tabwriter"

	"github.com/campbel/yoshi/options"
	"github.com/campbel/yoshi/types"
)

type Node struct {
	path     []string
	commands *types.OrderedMap[string, *Node]
	run      reflect.Value
}

func (n *Node) Exec(args ...string) error {
	if len(args) == 0 {
		return n.Run()
	}
	if node, ok := n.commands.Get(args[0]); ok {
		return node.Exec(args[1:]...)
	}
	return n.Run(args...)
}

func (n *Node) Run(args ...string) error {
	if !n.run.IsValid() {
		if len(args) == 0 {
			return fmt.Errorf("no command specified")
		}
		return fmt.Errorf("command not found: %s", args[0])
	}
	params, err := options.CreateFromArgs(n.options(), args)
	if err != nil {
		return err
	}
	retVals := n.run.Call(params)
	if len(retVals) > 0 {
		if err, ok := retVals[0].Interface().(error); ok {
			return err
		}
	}
	return nil
}

func (n *Node) options() []reflect.Type {
	if !n.run.IsValid() {
		return []reflect.Type{}
	}
	var params []reflect.Type
	paramCount := n.run.Type().NumIn()
	for i := 0; i < paramCount; i++ {
		params = append(params, n.run.Type().In(i))
	}
	return params
}

func (n *Node) Help() string {
	opts := n.options()

	var help string
	help += "Usage: " + strings.Join(n.path, " ")
	if n.commands.Len() > 0 {
		help += " COMMAND"
	}
	if len(opts) > 0 {
		nonPositional := false
		for _, opt := range options.GetOptions(opts[0]) {
			if opt.Positional() {
				help += " " + strings.ToUpper(opt.Flags[0])
			} else {
				nonPositional = true
			}
		}
		if nonPositional {
			help += " [OPTIONS]"
		}
	}
	if n.commands.Len() > 0 {
		help += "\nCommands:"
		n.commands.Each(func(key string, node *Node) {
			help += "\n  " + key
		})
	}
	if len(opts) > 0 {
		help += "\nOptions:"
		var buffer bytes.Buffer
		w := tabwriter.NewWriter(&buffer, 0, 0, 1, ' ', 0)
		for _, opt := range options.GetOptions(opts[0]) {
			line := "\n  " + strings.Join(opt.Flags, ", ")
			line += "\t" + opt.Type
			if opt.Description != "" {
				line += "\t" + fmt.Sprintf(`"%s"`, opt.Description)
			}
			if opt.Default != "" {
				line += fmt.Sprintf(` (default: "%s")`, opt.Default)
			}
			fmt.Fprint(w, line)
		}
		w.Flush()
		help += buffer.String()
	}
	return help
}

func NewTree(v any, path ...string) *Node {
	return parse(path, reflect.ValueOf(v))
}

func parse(path []string, v reflect.Value) *Node {
	switch v.Type().Kind() {
	case reflect.Struct:
		return parseStruct(path, v)
	case reflect.Func:
		return parseFunc(path, v)
	}
	return nil
}

func parseFunc(path []string, v reflect.Value) *Node {
	return &Node{
		path:     path,
		commands: types.NewOrderedMap[string, *Node](),
		run:      v,
	}
}

func parseStruct(path []string, v reflect.Value) *Node {
	n := &Node{
		path:     path,
		commands: types.NewOrderedMap[string, *Node](),
	}
	for _, field := range reflect.VisibleFields(v.Type()) {
		name := strings.ToLower(field.Name)
		n.commands.Set(name, parse(append(path, name), v.FieldByName(field.Name)))
	}
	return n
}

func (n *Node) Traverse(path ...string) *Node {
	if len(path) == 0 {
		return n
	}
	node, ok := n.commands.Get(path[0])
	if !ok {
		return nil
	}
	return node.Traverse(path[1:]...)
}

func (n *Node) TryTraverse(path ...string) (*Node, []string) {
	if len(path) == 0 {
		return n, nil
	}
	node, ok := n.commands.Get(path[0])
	if !ok {
		return n, path
	}
	return node.TryTraverse(path[1:]...)
}
