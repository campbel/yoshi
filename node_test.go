package yoshi

import (
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNodeBuild(t *testing.T) {
	t.Run("empty struct", func(t *testing.T) {
		type foo struct{}
		verifyNode(t, buildNodes(reflect.ValueOf(&foo{})),
			0, 0, reflect.ValueOf(foo{}).Type().String(), reflect.ValueOf(&foo{}).Type().String(), false, false)
	})
	t.Run("options only", func(t *testing.T) {
		type foo struct {
			Options struct{}
		}
		verifyNode(t, buildNodes(reflect.ValueOf(&foo{})),
			0, 0, reflect.ValueOf(foo{}).Type().String(), reflect.ValueOf(&foo{}).Type().String(), true, false)
	})
	t.Run("run only", func(t *testing.T) {
		type foo struct {
			Run func()
		}
		verifyNode(t, buildNodes(reflect.ValueOf(&foo{})),
			0, 0, reflect.ValueOf(foo{}).Type().String(), reflect.ValueOf(&foo{}).Type().String(), false, true)
	})
	t.Run("options and run", func(t *testing.T) {
		type foo struct {
			Options struct{}
			Run     func()
		}
		verifyNode(t, buildNodes(reflect.ValueOf(&foo{})),
			0, 0, reflect.ValueOf(foo{}).Type().String(), reflect.ValueOf(&foo{}).Type().String(), true, true)
	})
}

func verifyNode(t *testing.T, node *Node, commands, options int, tt, vt string, optsValid, runValid bool) {
	t.Helper()
	assert.NotNil(t, node)
	assert.Len(t, node.Commands, 0)
	assert.Len(t, node.Options, 0)
	assert.Equal(t, tt, node.Type.String())
	assert.Equal(t, vt, node.Value.Type().String())
	assert.Equal(t, optsValid, node.Opts.IsValid())
	assert.Equal(t, runValid, node.Run.IsValid())
}
