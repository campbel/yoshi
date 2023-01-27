package yoshi

import "testing"

type NodeApp struct {
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

func TestNodeImport(t *testing.T) {
	description := Describe[NodeApp]()
	if len(description.Commands) != 3 {
		t.Errorf("Expected 3 subcommands, got %d", len(description.Commands))
	}
}

func TestValidate(t *testing.T) {
	description := Describe[NodeApp]()
	if err := description.Validate(); err != nil {
		t.Errorf("Expected no errors, got %v", err)
	}
}
