package app

import (
	"testing"

	"github.com/campbel/yoshi/opts"
)

type RootOptions struct {
	Verbose bool   `opts:"-v,--verbose"`
	Output  string `opts:"-o,--output"`
}

type SomeOptions struct {
	Name string `opts:"-n,--name"`
}

type OtherOptions struct {
	Port int `opts:"-p,--port"`
}

func TestApp(t *testing.T) {
	var (
		rootOpts  RootOptions
		someOpts  SomeOptions
		otherOpts OtherOptions
	)

	app := App(func(args []string) {
		rootOpts = opts.MustParse[RootOptions](args)
	})
	app.Sub("some", func(args []string) {
		someOpts = opts.MustParse[SomeOptions](args)
	})
	app.Sub("other", func(args []string) {
		otherOpts = opts.MustParse[OtherOptions](args)
	})

	app.Parse([]string{"-v", "some", "-n", "foo"})
	if !rootOpts.Verbose {
		t.Error("rootOpts.Verbose should be true")
	}
	if someOpts.Name != "foo" {
		t.Error("someOpts.Name should be foo")
	}
	if otherOpts.Port != 0 {
		t.Error("otherOpts.Port should be 0")
	}

	app.Parse([]string{"other", "-p", "8080"})
	if otherOpts.Port != 8080 {
		t.Error("otherOpts.Port should be 8080")
	}
}
