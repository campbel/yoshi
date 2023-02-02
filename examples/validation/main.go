package main

import (
	"errors"
	"unicode"

	"github.com/campbel/yoshi"
)

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
