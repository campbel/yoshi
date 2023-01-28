package test

import (
	"bytes"
	"testing"

	"github.com/campbel/yoshi"
	"github.com/stretchr/testify/assert"
)

func TestYoshiSingleFunction(t *testing.T) {
	var buffer bytes.Buffer
	app := yoshi.New("test").WithConfig(yoshi.Config{HelpWriter: &buffer})

	type Options struct {
		Name string `yoshi:"-n,--name"`
	}
	var out Options

	app.Run(func(options Options) {
		out = options
	}, "-n", "mario")

	assert.Equal(t, "mario", out.Name)
}

func TestYoshiMultiFunction(t *testing.T) {

	type FetchOptions struct {
		URL     string `yoshi:"-u,--url"`
		Scheme  string `yoshi:"-s,--scheme" yoshi-default:"https"`
		Ignored string
	}
	type PrintOptions struct {
		Text      string `yoshi:"-t,--text"`
		Lowercase bool   `yoshi:"-l,--lowercase"`
	}
	type App struct {
		Fetch func(options FetchOptions)
		Print func(options PrintOptions)
	}

	t.Run("fetch", func(t *testing.T) {
		var out FetchOptions
		var buffer bytes.Buffer
		yoshi.New("test").WithConfig(yoshi.Config{HelpWriter: &buffer}).
			Run(App{
				Fetch: func(options FetchOptions) {
					out = options
				},
				Print: func(options PrintOptions) {
					t.Fatal("Print should not be called")
				},
			}, "fetch", "-u", "https://google.com")
		assert.Equal(t, "https://google.com", out.URL)
		assert.Equal(t, "https", out.Scheme)
		assert.Equal(t, "", buffer.String())
	})

	t.Run("print", func(t *testing.T) {
		var out PrintOptions
		var buffer bytes.Buffer
		yoshi.New("test").WithConfig(yoshi.Config{HelpWriter: &buffer}).
			Run(App{
				Fetch: func(options FetchOptions) {
					t.Fatal("Fetch should not be called")
				},
				Print: func(options PrintOptions) {
					out = options
				},
			}, "print", "-t", "what is going on?")
		assert.Equal(t, "what is going on?", out.Text)
		assert.Equal(t, "", buffer.String())
	})

	t.Run("print lowercase", func(t *testing.T) {
		var out PrintOptions
		var buffer bytes.Buffer
		app := App{
			Fetch: func(options FetchOptions) {
				t.Fatal("Fetch should not be called")
			},
			Print: func(options PrintOptions) {
				out = options
			},
		}
		yosh := yoshi.New("test").WithConfig(yoshi.Config{HelpWriter: &buffer})

		yosh.Run(app, "print", "-l", "-t", "what is going on?")
		assert.Equal(t, "what is going on?", out.Text)
		assert.Equal(t, true, out.Lowercase)
		assert.Equal(t, "", buffer.String())

		out = PrintOptions{}
		yosh.Run(app, "print", "-l", "true", "-t", "who are you?")
		assert.Equal(t, "who are you?", out.Text)
		assert.Equal(t, true, out.Lowercase)
		assert.Equal(t, "", buffer.String())
	})

	t.Run("funch", func(t *testing.T) {
		var out FetchOptions
		var buffer bytes.Buffer
		yoshi.New("test").WithConfig(yoshi.Config{HelpWriter: &buffer}).
			Run(App{
				Fetch: func(options FetchOptions) {
					out = options
				},
				Print: func(options PrintOptions) {
					t.Fatal("Print should not be called")
				},
			}, "funch", "-u", "https://google.com")
		assert.Empty(t, out.URL)
		assert.Equal(t, "Error: command not found: funch\nUsage: test COMMAND\nCommands:\n  fetch\n  print\n", buffer.String())
	})

	t.Run("bad flag", func(t *testing.T) {
		var out FetchOptions
		var buffer bytes.Buffer
		yoshi.New("test").WithConfig(yoshi.Config{HelpWriter: &buffer}).
			Run(App{
				Fetch: func(options FetchOptions) {
					out = options
				},
				Print: func(options PrintOptions) {
					t.Fatal("Print should not be called")
				},
			}, "fetch", "-t", "https://google.com")
		assert.Empty(t, out.URL)
		assert.Equal(t, "Error: unknown flag: -t\nUsage: test fetch [OPTIONS]\nOptions:\n  -u,--url    string\n  -s,--scheme string (default: https)\n", buffer.String())
	})
}
