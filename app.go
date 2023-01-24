package yoshi

import (
	"os"
)

type Command struct {
	Name string
	Run  []func([]string)
	subs []*Command
}

func App(run ...func([]string)) *Command {
	return &Command{Name: "root", Run: run}
}

func (n *Command) Sub(name string, fn ...func([]string)) *Command {
	n.subs = append(n.subs, &Command{
		Name: name,
		Run:  fn,
	})
	return n
}

func (n *Command) Start() {
	n.Parse(os.Args[1:])
}

func (n *Command) Parse(args []string) {
	if n == nil {
		return
	}
	sub, i := getFirstSub(n.subs, args)
	for _, fn := range n.Run {
		fn(args[:i])
	}
	if sub != nil {
		sub.Parse(args[i+1:])
	}
}

func getFirstSub(subs []*Command, args []string) (*Command, int) {
	for i, arg := range args {
		for _, sub := range subs {
			if sub.Name == arg {
				return sub, i
			}
		}
	}
	return nil, len(args)
}
