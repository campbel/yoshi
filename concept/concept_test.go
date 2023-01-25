package concept

import (
	"strings"
	"testing"
)

type NuApp struct {
	Options struct {
		Verbose bool `yoshi:"flag=-v,--verbose"`
	}
	Server struct {
		Options ServerOptions
		Run     func(ServerOptions)
	}
	Client struct {
		Options ClientOptions
		Run     func(ClientOptions)
	}
}

type ClientOptions struct {
	Name     string            `yoshi:"flag=-n,--name"`
	Insecure bool              `yoshi:"flag=-k,--insecure"`
	Count    int               `yoshi:"flag=-c,--count"`
	ListStr  []string          `yoshi:"flag=-l,--list"`
	ListInt  []int             `yoshi:"flag=-i,--int-list"`
	ListBool []bool            `yoshi:"flag=-b,--bool-list"`
	DictStr  map[string]string `yoshi:"flag=-d,--dict"`
	DictInt  map[string]int    `yoshi:"flag=-e,--int-dict"`
	DictBool map[string]bool   `yoshi:"flag=-f,--bool-dict"`
}

type ServerOptions struct {
	Port string `yoshi:"flag=-p,--port"`
}

func TestConcept(t *testing.T) {
	count := 0
	app := new(NuApp)
	app.Client.Run = func(options ClientOptions) {
		if options.Name != "bob" {
			t.Error("name should be bob")
		}
		count++
	}
	app.Server.Run = func(options ServerOptions) {
		t.Fatal("server should not run")
	}

	args := "client --name bob --insecure true --count 1 "
	args += "--list 1,2,3 --int-list 1,2,3 --bool-list true,false,true "
	args += "--dict a=d,b=e,c=f --int-dict a=1,b=2,c=3 --bool-dict a=true,b=false,c=true"
	Run(app, strings.Split(args, " ")...)
	if count != 1 {
		t.Error("count should be 2")
	}
}
