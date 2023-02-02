package main

import (
	"fmt"

	"cuelang.org/go/cue/cuecontext"
	"github.com/campbel/yoshi"
)

// Use a validation schema rather than cue struct tags, as its more powerful
const optionsSchema = `
First: !=""
Middle?: string
Last: !=""
FullName: [
  if Middle == "" { First + " " + Last }
  if Middle != "" { First + " " + Middle + " " + Last }
][0]
`

type Options struct {
	First    string `yoshi:"-f,--first;Name of the person, required and must be capitalized;"`
	Middle   string `yoshi:"-m,--middle;Name of the person, options but must be capitalized;"`
	Last     string `yoshi:"-l,--last;Name of the person, required and must be capitalized;"`
	FullName string `json:",omitempty"`
}

func main() {
	yoshi.New("test").Run(cueValidate(optionsSchema, func(options Options) error {
		fmt.Printf("Hello %v", options.FullName)
		return nil
	}))
}

// A generic cue validation middleware
func cueValidate[T any](schema string, next func(T) error) func(T) error {
	return func(options T) error {
		ctx := cuecontext.New()
		val := ctx.CompileString(schema).Unify(ctx.Encode(options))
		if err := val.Err(); err != nil {
			return err
		}
		if err := val.Validate(); err != nil {
			return err
		}
		// Decode the value back into options, if you want completed values (ex: FullName)
		if err := val.Decode(&options); err != nil {
			return err
		}
		return next(options)
	}
}
