package yoshi

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestApp(t *testing.T) {
	assert := assert.New(t)
	ctx := Create[TestYoshiApp]("appy")
	if err := ctx.Validate(); len(err) != 4 {
		t.Error(err)
	}
	assert.Equal("Usage: appy call message [options]\nOptions:\n  -t string The text to send\n", ctx.help("call", "-n", "123", "message", "-t", "hello, world"))
	ctx.run("call", "-n", "123", "message", "-t", "hello, world")
	if ctx.App.Call.callCount != 1 {
		t.Errorf("Expected 1, got %d", ctx.App.Call.callCount)
	}

	address := ""
	ctx.App.Email.Run = func(options struct {
		Address string `yoshi-flag:"-a"`
	}) {
		address = options.Address
	}
	ctx.run("email", "-a", "foo@bar")
	if address != "foo@bar" {
		t.Errorf("Expected foo@bar, got %s", address)
	}

	ctx.App.Text.Run = func(options TextOptions) {
		address = options.Number
	}
	ctx.run("text", "-n", "123")
	if address != "123" {
		t.Errorf("Expected 123, got %s", address)
	}
	if ctx.App.Options.Name != "smudge" {
		t.Errorf("Expected smudge, got %s", ctx.App.Options.Name)
	}
}

type TestYoshiApp struct {
	Options struct {
		Name    string `yoshi-flag:"-n" yoshi-def:"smudge"`
		Address string `yoshi-flag:"-a"`
	}
	Call  CallCommand
	Email struct {
		Options struct {
			Address string `yoshi-flag:"-a"`
		}
		Run func(options struct {
			Address string `yoshi-flag:"-a"`
		})
	}
	Text struct {
		Options TextOptions
		Run     func(options TextOptions)
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

type TextOptions struct {
	Number string `yoshi-flag:"-n"`
}
