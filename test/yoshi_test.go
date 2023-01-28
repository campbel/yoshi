package test

import (
	"testing"

	"github.com/campbel/yoshi"
	"github.com/stretchr/testify/assert"
)

func TestEmptyApp(t *testing.T) {
	ctx := yoshi.Create[struct{}]("empty")
	if err := ctx.Validate(); len(err) != 0 {
		t.Error(err)
	}
}

func TestSimpleApp(t *testing.T) {
	type SimpleApp struct {
		Options struct {
			Name string `yoshi-flag:"-n" yoshi-desc:"The name of the user"`
		}
	}
	ctx := yoshi.Create[SimpleApp]("simple")
	if err := ctx.Validate(); len(err) != 0 {
		t.Error(err)
	}
	ctx.Run("simple", "-n", "smudge")
	assert.Equal(t, "smudge", ctx.App.Options.Name)
	assert.Equal(t, "Usage: simple [options]\nOptions:\n  -n string The name of the user\n", ctx.Help("simple", "-n", "smudge"))
}

func TestSimpleAppWithDefault(t *testing.T) {
	type SimpleApp struct {
		Options struct {
			Name string `yoshi-flag:"-n" yoshi-desc:"The name of the user" yoshi-def:"smidge"`
		}
	}
	ctx := yoshi.Create[SimpleApp]("simple")
	if err := ctx.Validate(); len(err) != 0 {
		t.Error(err)
	}
	ctx.Run("simple")
	assert.Equal(t, "smidge", ctx.App.Options.Name)
	assert.Equal(t, "Usage: simple [options]\nOptions:\n  -n string The name of the user (default: smidge)\n", ctx.Help("simple"))
}

func TestSimpleAppWithRun(t *testing.T) {
	type Options struct {
		Name string `yoshi-flag:"-n" yoshi-desc:"The name of the user"`
	}
	type SimpleApp struct {
		Options Options
		Run     func(Options)
	}
	ctx := yoshi.Create[SimpleApp]("simple")
	if err := ctx.Validate(); len(err) != 0 {
		t.Error(err)
	}
	var called bool
	ctx.App.Run = func(options Options) {
		called = true
	}
	ctx.Run("simple", "-n", "smudge")
	assert.True(t, called)
}

func TestSimpleAppWithBadOptionsType(t *testing.T) {
	type SimpleApp struct {
		Options string
	}
	ctx := yoshi.Create[SimpleApp]("simple")
	if err := ctx.Validate(); len(err) != 1 {
		t.Error(err)
	}
}
