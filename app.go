package yoshi

import (
	"os"
)

type Command struct {
	Name string
	Run  []Runner
	subs []*Command
}

func App(run ...Runner) *Command {
	return &Command{Name: "root", Run: run}
}

func (n *Command) Sub(name string, fn ...Runner) *Command {
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
	for _, runner := range n.Run {
		runner.Run(args[:i])
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

type Args []string

func (a Args) Help() bool {
	for _, arg := range a {
		if arg == "--help" {
			return true
		}
	}
	return false
}

type Runner interface {
	Run(Args)
}

type RunMulti struct {
	fns []RunnerFunc
}

func (r RunMulti) Run(args Args) {
	for _, fn := range r.fns {
		fn.Run(args)
	}
}

func Run(fns ...RunnerFunc) RunMulti {
	return RunMulti{fns}
}

type RunnerFunc func(Args)

func (r RunnerFunc) Run(args []string) {
	r(args)
}
