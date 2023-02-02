# Yoshi

Yoshi is a minimalist framework for command line applications. The goal of Yoshi is to minimize boilerplate and setup while maintaining modern CLI application standards, such as auto-generated help text, default values, positional arguments and short and long flags.

Yoshi uses struct tags to define flag names, description and default values. The format of the tag is:

```golang
yoshi:"[FLAG1],[FLAG2];[DESCRIPTION];[DEFAULT VALUE]"
```

Tags should be added to types passed to functions in your app struct.

## Examples

Learning from examples is great, hopefully this is sufficient inplace of actual docs.

### [Basic](/examples/basic/main.go)

A minimal example application takes a function to run with the options structure to parse.

Everything is optional, but missing the flags means Yoshi will ignore the field.

```golang
package main

import "github.com/campbel/yoshi"

type Options struct {
  Message string `yoshi:"MESSAGE;The message to print;Hello, world!"`
}

func main() {
  yoshi.New("basic").Run(func(opts Options) {
    println(opts.Message)
  })
}
```

Default values are used, even for positional parameters.

```bash
go run main.go
Hello, world!
```

They can be passed as expected.

```bash
go run main.go "This is my message"
This is my message
```

Help text, automatically generated.

```bash
go run main.go --help
Usage: basic MESSAGE
Options:
  MESSAGE string "The message to print" (default: "Hello, world!")
```

### [Typical](/examples/typical/main.go)

A typical application will likely have multiple commands. This is managed through structs, like so.

```golang
type FetchOptions struct {
  URL     string            `yoshi:"URL;The URL to fetch;"`
  Method  string            `yoshi:"-m,--method;The HTTP method to use;GET"`
  Body    string            `yoshi:"-b,--body;The request body;"`
  Headers map[string]string `yoshi:"-H,--header;The request headers;"`
}

type ServeOptions struct {
  Dir  string `yoshi:"DIRECTORY;The directory to serve files from;."`
  Port int    `yoshi:"-p,--port;The port to serve on;8080"`
}

type App struct {
  Fetch func(FetchOptions)
  Serve func(ServeOptions)
}


func main() {
  yoshi.New("typical").Run(App{
    Fetch: func(opts FetchOptions) {
      // TODO
    },
    Serve: func(opts ServeOptions) {
      // TODO
    },
  })
}
```

Useful help text still.

```bash
go run main.go fetch --help
Usage: typical fetch URL [OPTIONS]
Options:
  URL          string            "The URL to fetch"
  -m, --method string            "The HTTP method to use" (default: "GET")
  -b, --body   string            "The request body"
  -h, --header map[string]string "The request headers"
```

Usage is as you'd expect

```bash
go run main.go fetch -m POST -b '{"foo":"bar"}' -H "Content-Type=application/json" http://httpbin.org/post
```

### [Validation](/examples/validation/main.go)

You may want additional validation beyond what Yoshi enforces for you. To do this, perform the validation as a first step of the function execution and return any errors.

```golang
type Options struct {
  Name string `yoshi:"-n,--name;Name of the person, required and must be capitalized;"`
}

func main() {
  yoshi.New("test").Run(func(options Options) error {
    if options.Name == "" {
      return errors.New("name is required")
    }
    if !unicode.IsUpper(rune(options.Name[0])) {
      return errors.New("name must be capitalized")
    }
    return nil
  })
}
```

If an error is returned, Yoshi will print it along with the help text.

```bash
go run main.go
error: name is required
Usage: test [OPTIONS]
Options:
  -n, --name string "Name of the person, required and must be capitalized"
```

### [Yoshi <3 Cue](/examples/cue/main.go)

There is no direct integration between Yoshi and Cue, but since we're working with Go structs, there is seamless cooperation between the two.

1. Setup a simple Yoshi handler with options.

    ```golang
    type Options struct {
      Value `yoshi:"-v"`
    }

    func main() {
      yoshi.New("cue-example").Run(func(options Options) {
        // implemented in step 3 below
      })
    }
    ```

1. Next, create a cue schema for the options type and load it as a cue value.

    ```golang
    const schema = `
    Value != ""
    `
    ```

1. Last, perform a cue validation as a first step in the handler

    ```golang
      yoshi.New("cue-example").Run(func(options Options) {
        ctx := cuecontext.New()
        val := ctx.CompileString(optionsSchema).Unify(ctx.Encode(options))
        if err := val.Err(); err != nil {
          return err
        }
        if err := val.Validate(); err != nil {
          return err
        }
        // TODO implement functionality
      })
    ```

This is the basic idea, but check out the example code for a more realistic implementation including a generic validation wrapper.
