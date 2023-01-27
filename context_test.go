package yoshi

import (
	"fmt"
	"testing"
)

func TestApp(t *testing.T) {
	ctx := Create[TestYoshiApp]("appy")
	if err := ctx.Validate(); len(err) != 4 {
		t.Error(err)
	}
	helpText := ctx.help("call", "-n", "123", "message", "-t", "hello, world")
	fmt.Print(helpText)
	ctx.run("call", "-n", "123", "message", "-t", "hello, world")

	if ctx.App.Call.callCount != 1 {
		t.Errorf("Expected 1, got %d", ctx.App.Call.callCount)
	}
}

type TestYoshiApp struct {
	Options struct {
		Name    string `yoshi-flag:"-n"`
		Address string `yoshi-flag:"-a"`
	}
	Call  CallCommand
	Email struct {
		Options struct {
			Address string `yoshi-flag:"-a"`
		}
	}
	Text struct {
		Options struct {
			Number string `yoshi-flag:"-n"`
		}
	}
}

type CallCommand struct {
	callCount int
	Options   CallOptions
	Message   struct {
		Options struct {
			Text string `yoshi-flag:"-t" yoshi-desc:"The text to send"`
		}
	}
}

func (c *CallCommand) Run(options CallOptions) {
	c.callCount++
}

type CallOptions struct {
	Number int `yoshi-flag:"-n" yoshi-desc:"The number to call"`
}
