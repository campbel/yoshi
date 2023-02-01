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
