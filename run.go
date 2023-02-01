package yoshi

import (
	"fmt"
	"io"
	"os"

	"github.com/campbel/yoshi/parser"
)

type Config struct {
	HelpWriter io.Writer
}

var defaultConfig = Config{HelpWriter: os.Stdout}

type Yoshi struct {
	name   string
	config Config
}

func New(name string) *Yoshi {
	return &Yoshi{name: name, config: defaultConfig}
}

func (y *Yoshi) WithConfig(config Config) *Yoshi {
	y.config = config
	return y
}

func (y *Yoshi) Run(v any, args ...string) error {
	if len(args) == 0 {
		args = os.Args[1:]
	}
	return y.RunWithArgs(v, args...)
}

func (y *Yoshi) RunWithArgs(v any, args ...string) error {
	root := parser.NewTree(v, y.name)
	node, leftOver := root.TryTraverse(args...)
	if hasHelp(leftOver) {
		fmt.Fprintln(y.config.HelpWriter, node.Help())
		return nil
	}
	err := node.Run(leftOver...)
	if err != nil {
		fmt.Fprintln(y.config.HelpWriter, "error: "+err.Error())
		fmt.Fprintln(y.config.HelpWriter, node.Help())
		return err
	}
	return nil
}

func hasHelp(args []string) bool {
	for _, arg := range args {
		if arg == "-h" || arg == "--help" {
			return true
		}
	}
	return false
}
