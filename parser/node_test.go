package parser

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewNode(t *testing.T) {
	type FetchOptions struct {
		URL    string `yoshi:";URL to fetch;"`
		Method string `yoshi:"-m,--method;HTTP request method;GET"`
	}

	type EchoOptions struct {
		Message string `yoshi:";Message to echo;"`
	}

	type App struct {
		Foo struct {
			Fetch func(FetchOptions)
		}
		Bar struct {
			Echo func(EchoOptions)
		}
		Baz struct {
			Print func()
		}
	}

	fetchChan := make(chan (FetchOptions))
	echoChan := make(chan (EchoOptions))
	app := App{}
	app.Foo.Fetch = func(options FetchOptions) {
		go func() {
			fetchChan <- options
		}()
	}
	app.Bar.Echo = func(options EchoOptions) {
		go func() {
			echoChan <- options
		}()
	}
	rootNode := NewTree(app)

	t.Run("node structure", func(t *testing.T) {
		assert := assert.New(t)
		assert.Equal(3, rootNode.commands.Len())
		assert.Equal(1, rootNode.commands.GetVal("foo").commands.Len())
		assert.Equal(1, rootNode.commands.GetVal("bar").commands.Len())
		assert.Equal(1, rootNode.commands.GetVal("baz").commands.Len())
		assert.Equal(rootNode.commands.GetVal("foo"), rootNode.Traverse("foo"))
		assert.Equal(rootNode.commands.GetVal("bar"), rootNode.Traverse("bar"))
		assert.Equal(rootNode.commands.GetVal("baz"), rootNode.Traverse("baz"))
		assert.Equal(rootNode.commands.GetVal("foo").commands.GetVal("fetch"), rootNode.Traverse("foo", "fetch"))
		assert.Equal(rootNode.commands.GetVal("bar").commands.GetVal("echo"), rootNode.Traverse("bar", "echo"))
		assert.Equal(rootNode.commands.GetVal("baz").commands.GetVal("print"), rootNode.Traverse("baz", "print"))
		assert.Nil(rootNode.Traverse("foo", "fetch", "echo"))
	})

	t.Run("run fetch", func(t *testing.T) {
		assert := assert.New(t)
		if err := rootNode.Traverse("foo", "fetch").Run(); err != nil {
			t.Fatal(err)
		}
		fetchOptions := <-fetchChan
		assert.Equal(fetchOptions.Method, "GET")
		assert.Empty(fetchOptions.URL)

		if err := rootNode.Traverse("foo", "fetch").Run("-m", "POST", "http://example.com"); err != nil {
			t.Fatal(err)
		}
		fetchOptions = <-fetchChan
		assert.Equal(fetchOptions.Method, "POST")
		assert.Equal(fetchOptions.URL, "http://example.com")

		if err := rootNode.Exec("foo", "fetch", "-m", "POST", "http://example.com"); err != nil {
			t.Fatal(err)
		}
		fetchOptions = <-fetchChan
		assert.Equal(fetchOptions.Method, "POST")
		assert.Equal(fetchOptions.URL, "http://example.com")

		if err := rootNode.Traverse("bar", "echo").Run("hello"); err != nil {
			t.Fatal(err)
		}
	})

	t.Run("exec echo", func(t *testing.T) {
		assert := assert.New(t)
		echoOptions := <-echoChan
		assert.Equal(echoOptions.Message, "hello")

		if err := rootNode.Exec("bar", "echo", "hello"); err != nil {
			t.Fatal(err)
		}

		echoOptions = <-echoChan
		assert.Equal(echoOptions.Message, "hello")
	})

	t.Run("fetch help", func(t *testing.T) {
		assert := assert.New(t)
		fetchNode := rootNode.Traverse("foo", "fetch")
		expected := `Usage: foo fetch URL [OPTIONS]
Options:
  URL          string "URL to fetch"
  -m, --method string "HTTP request method" (default: "GET")`
		assert.Equal(expected, fetchNode.Help())
	})

	t.Run("echo help", func(t *testing.T) {
		assert := assert.New(t)
		echoNode := rootNode.Traverse("bar", "echo")
		expected := `Usage: bar echo MESSAGE
Options:
  MESSAGE string "Message to echo"`
		assert.Equal(expected, echoNode.Help())
	})

}
