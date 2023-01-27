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
	t.Run("sub command", func(t *testing.T) {
		type sub struct {
		}
		type foo struct {
			Sub sub
		}
		node := buildNodes(reflect.ValueOf(&foo{}))
		verifyNode(t, node,
			1, 0, reflect.ValueOf(foo{}).Type().String(), reflect.ValueOf(&foo{}).Type().String(), false, false)
		verifyNode(t, node.Commands["sub"],
			0, 0, reflect.ValueOf(sub{}).Type().String(), reflect.ValueOf(&sub{}).Type().String(), false, false)
	})
	t.Run("multiple sub commands", func(t *testing.T) {
		type sub1 struct {
		}
		type sub2 struct {
		}
		type foo struct {
			Sub1 sub1
			Sub2 sub2
		}
		node := buildNodes(reflect.ValueOf(&foo{}))
		verifyNode(t, node,
			2, 0, reflect.ValueOf(foo{}).Type().String(), reflect.ValueOf(&foo{}).Type().String(), false, false)
		verifyNode(t, node.Commands["sub1"],
			0, 0, reflect.ValueOf(sub1{}).Type().String(), reflect.ValueOf(&sub1{}).Type().String(), false, false)
		verifyNode(t, node.Commands["sub2"],
			0, 0, reflect.ValueOf(sub2{}).Type().String(), reflect.ValueOf(&sub2{}).Type().String(), false, false)
	})
}

func verifyNode(t *testing.T, node *Node, commandCount, optionsCount int, tt, vt string, optsValid, runValid bool) {
	t.Helper()
	assert.NotNil(t, node)
	assert.Len(t, node.Commands, commandCount)
	assert.Len(t, node.Options, optionsCount)
	assert.Equal(t, tt, node.Type.String())
	assert.Equal(t, vt, node.Value.Type().String())
	assert.Equal(t, optsValid, node.Opts.IsValid())
	assert.Equal(t, runValid, node.Run.IsValid())
}
