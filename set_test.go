package yoshi

import (
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSetter(t *testing.T) {
	t.Run("bool", func(t *testing.T) {
		assert := assert.New(t)
		assert.NoError(runSetterTest(t, "true", true))
		assert.NoError(runSetterTest(t, "false", false))
		assert.NoError(runSetterTest(t, "", true))
		assert.Error(runSetterTest(t, "qweqweqwe", false))
	})
	t.Run("string", func(t *testing.T) {
		assert := assert.New(t)
		assert.NoError(runSetterTest(t, "", ""))
		assert.NoError(runSetterTest(t, "blugh", "blugh"))
	})
	t.Run("int", func(t *testing.T) {
		assert := assert.New(t)
		assert.NoError(runSetterTest(t, "32", 32))
		assert.NoError(runSetterTest(t, "-32", -32))
		assert.Error(runSetterTest(t, "", 0))
		assert.Error(runSetterTest(t, "toad", 0))
	})
	t.Run("bool slice", func(t *testing.T) {
		assert := assert.New(t)
		assert.NoError(runSetterTest(t, "true,false,true", []bool{true, false, true}))
		assert.NoError(runSetterTest(t, "", []bool{}))
		assert.Error(runSetterTest(t, "asdads", []bool{}))
	})
	t.Run("string slice", func(t *testing.T) {
		assert := assert.New(t)
		assert.NoError(runSetterTest(t, "true,false,true", []string{"true", "false", "true"}))
		assert.NoError(runSetterTest(t, "", []string{}))
		assert.NoError(runSetterTest(t, "asdads", []string{"asdads"}))
	})
	t.Run("int slice", func(t *testing.T) {
		assert := assert.New(t)
		assert.NoError(runSetterTest(t, "1,2,3", []int{1, 2, 3}))
		assert.Error(runSetterTest(t, "1,", []int{1}))
		assert.NoError(runSetterTest(t, "", []int{}))
		assert.Error(runSetterTest(t, "asdads", []int{}))
	})
	t.Run("bool map", func(t *testing.T) {
		assert := assert.New(t)
		assert.NoError(runSetterTest(t, "foo=true,bar=false,baz=true", map[string]bool{"foo": true, "bar": false, "baz": true}))
		assert.NoError(runSetterTest(t, "", map[string]bool{}))
		assert.Error(runSetterTest(t, "foo", map[string]bool{}))
		assert.NoError(runSetterTest(t, "foo=true", map[string]bool{"foo": true}))
		assert.Error(runSetterTest(t, "foo=qweqwe", map[string]bool{}))
	})
	t.Run("string map", func(t *testing.T) {
		assert := assert.New(t)
		assert.NoError(runSetterTest(t, "foo=bar,fast=slow", map[string]string{"foo": "bar", "fast": "slow"}))
		assert.NoError(runSetterTest(t, "", map[string]string{}))
		assert.Error(runSetterTest(t, "asdads", map[string]string{}))
	})
	t.Run("int map", func(t *testing.T) {
		assert := assert.New(t)
		assert.NoError(runSetterTest(t, "foo=1,bar=2,baz=3", map[string]int{"foo": 1, "bar": 2, "baz": 3}))
		assert.NoError(runSetterTest(t, "", map[string]int{}))
		assert.Error(runSetterTest(t, "asdasd", map[string]int{}))
		assert.Error(runSetterTest(t, "foo=qweqwe", map[string]int{}))
	})
	t.Run("bad map", func(t *testing.T) {
		assert := assert.New(t)
		assert.Error(runSetterTest(t, "true=1,false=2", map[bool]int(nil)))
	})
}

func runSetterTest[T any](t *testing.T, input string, expected T) error {
	var v T
	value := reflect.ValueOf(&v)
	setter, ok := setterMap[value.Elem().Kind()]
	assert.True(t, ok)
	err := setter(value.Elem(), input)
	assert.Equal(t, expected, v)
	return err
}
