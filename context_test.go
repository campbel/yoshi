package yoshi

import "testing"

func TestApp(t *testing.T) {
	count := 0
	ctx := Create[TestYoshiApp]()
	if err := ctx.Validate(); err != nil {
		t.Error(err)
	}
	ctx.App.Call.Run = func() {
		count++
	}
	ctx.run("call", "-n", "123", "message", "-t", "hello, world")

	if count != 1 {
		t.Errorf("Expected 1, got %d", count)
	}
}

type TestYoshiApp struct {
	Options struct {
		Name    string `yoshi-flag:"-n"`
		Address string `yoshi-flag:"-a"`
	}
	Call struct {
		Options struct {
			Number int `yoshi-flag:"-n" yoshi-def:"hello, world"`
		}
		Message struct {
			Options struct {
				Text string `yoshi-flag:"-t"`
			}
		}
		Run func()
	}
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
