package yoshi

import (
	"bytes"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

type FuncOptions struct {
	Name string `yoshi:"-n,--name"`
}

func TestYoshiFunc(t *testing.T) {
	Run(func(options FuncOptions) {
		fmt.Println(options)
	}, "-n", "mario")
}

type PrintOptions struct {
	Text string `yoshi:"-t,--text"`
}

type FetchOptions struct {
	URL string `yoshi:"-u,--url"`
}

type testApp struct {
	Print func(options PrintOptions)
	Fetch func(options FetchOptions)
}

func TestYoshiObj(t *testing.T) {
	var buffer bytes.Buffer
	err := New("test").WithConfig(Config{&buffer}).Run(testApp{
		Print: func(options PrintOptions) {
			fmt.Println(options)
		},
		Fetch: func(options FetchOptions) {
			fmt.Println(options)
		},
	}, "funch", "-t", "http://google.com")
	if err == nil {
		assert.Fail(t, err.Error())
	}
	assert.Equal(t, "Error: command not found: funch\nUsage: test COMMAND\nCommands:\n  print\n  fetch\n", buffer.String())
}
