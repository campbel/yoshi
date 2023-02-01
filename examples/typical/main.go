package main

import (
	"bytes"
	"fmt"
	"io"
	"net/http"

	"github.com/campbel/yoshi"
)

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
			request, err := http.NewRequest(opts.Method, opts.URL, bytes.NewBufferString(opts.Body))
			if err != nil {
				panic(err)
			}
			for key, value := range opts.Headers {
				request.Header.Add(key, value)
			}
			response, err := http.DefaultClient.Do(request)
			if err != nil {
				panic(err)
			}
			data, _ := io.ReadAll(response.Body)
			fmt.Printf(`Status: %s
Headers: %v
Body: %s`, response.Status, response.Header, data)
		},
		Serve: func(opts ServeOptions) {
			fmt.Println("Serving files from", opts.Dir, "on port", opts.Port)
			err := http.ListenAndServe(fmt.Sprintf(":%d", opts.Port), http.FileServer(http.Dir(opts.Dir)))
			if err != nil {
				panic(err)
			}
		},
	})
}
