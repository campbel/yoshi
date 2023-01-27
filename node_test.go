package yoshi

import (
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNodeBuild(t *testing.T) {
	t.Run("empty struct", func(t *testing.T) {
		type foo struct{}
		verifyNode(t, buildNodes("root", reflect.ValueOf(&foo{})),
			0, 0, reflect.ValueOf(&foo{}).Type().String(), false, false)
	})
	t.Run("options only", func(t *testing.T) {
		type foo struct {
			Options struct{}
		}
		verifyNode(t, buildNodes("root", reflect.ValueOf(&foo{})),
			0, 0, reflect.ValueOf(&foo{}).Type().String(), true, false)
	})
	t.Run("run only", func(t *testing.T) {
		type foo struct {
			Run func()
		}
		verifyNode(t, buildNodes("root", reflect.ValueOf(&foo{})),
			0, 0, reflect.ValueOf(&foo{}).Type().String(), false, true)
	})
	t.Run("options and run", func(t *testing.T) {
		type foo struct {
			Options struct{}
			Run     func()
		}
		verifyNode(t, buildNodes("root", reflect.ValueOf(&foo{})),
			0, 0, reflect.ValueOf(&foo{}).Type().String(), true, true)
	})
	t.Run("sub command", func(t *testing.T) {
		type sub struct {
		}
		type foo struct {
			Sub sub
		}
		node := buildNodes("root", reflect.ValueOf(&foo{}))
		verifyNode(t, node,
			1, 0, reflect.ValueOf(&foo{}).Type().String(), false, false)
		verifyNode(t, node.commands["sub"],
			0, 0, reflect.ValueOf(&sub{}).Type().String(), false, false)
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
		node := buildNodes("root", reflect.ValueOf(&foo{}))
		verifyNode(t, node,
			2, 0, reflect.ValueOf(&foo{}).Type().String(), false, false)
		verifyNode(t, node.commands["sub1"],
			0, 0, reflect.ValueOf(&sub1{}).Type().String(), false, false)
		verifyNode(t, node.commands["sub2"],
			0, 0, reflect.ValueOf(&sub2{}).Type().String(), false, false)
	})
}

func verifyNode(t *testing.T, node *cmdNode, commandCount, optionsCount int, valueType string, optsValid, runValid bool) {
	t.Helper()
	assert.NotNil(t, node)
	assert.Len(t, node.commands, commandCount)
	assert.Len(t, node.options, optionsCount)
	assert.Equal(t, valueType, node.cmdReference.Type().String())
	assert.Equal(t, optsValid, node.optionsReference.IsValid())
	assert.Equal(t, runValid, node.runReference.IsValid())
}
